package screen

import (
	"errors"
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/vt"
)

// TestInstanceNew: New takes cols 1..2000 and rows 1..1000 (one column
// reads as vt.MinimumCols) and refuses anything else; the payload cap is
// Config's, DefaultPayloadLimit by default.
func TestInstanceNew(t *testing.T) {
	for _, s := range [][2]int{{0, 24}, {80, 0}, {-1, 24}, {2001, 24}, {80, 1001}} {
		if i, err := New(s[0], s[1], Config{}); err == nil || i != nil {
			t.Errorf("New(%d, %d) accepted", s[0], s[1])
		}
	}
	for _, s := range [][2]int{{1, 1}, {2000, 1000}, {80, 24}} {
		i, err := New(s[0], s[1], Config{})
		if err != nil {
			t.Fatalf("New(%d, %d): %v", s[0], s[1], err)
		}
		if want := max(s[0], vt.MinimumCols); i.t.Cols() != want || i.t.Rows() != s[1] {
			t.Errorf("New(%d, %d) is %dx%d", s[0], s[1], i.t.Cols(), i.t.Rows())
		}
	}
	if _, err := New(80, 24, Config{Unicode: "11"}); err == nil {
		t.Errorf("Unicode 11 accepted")
	}
	i, _ := New(80, 24, Config{})
	j, _ := New(80, 24, Config{PayloadLimit: 5})
	if i.t.Options.PayloadLimit != DefaultPayloadLimit || j.t.Options.PayloadLimit != 5 {
		t.Errorf("payload limits %d, %d", i.t.Options.PayloadLimit, j.t.Options.PayloadLimit)
	}
}

// TestInstanceWriteRead: Write takes the bytes, Read reads them, Resize
// refuses a size outside the limits and leaves the size, and re-caps the
// history for the new width as Resize does; Delta is the instance's cache.
func TestInstanceWriteRead(t *testing.T) {
	i, _ := New(100, 5, Config{})
	if n, err := i.Write([]byte("hello\x1b[1mbold")); n != 13 || err != nil {
		t.Fatalf("Write = %d, %v", n, err)
	}
	f, err := i.Read(Request{Attrs: true})
	if err != nil || f.Trimmed()[0] != "hellobold" || len(f.AttrRuns) != 1 {
		t.Fatalf("Read = %q, %d runs, %v", f.Trimmed()[0], len(f.AttrRuns), err)
	}
	if _, delta, _ := i.Delta(&f, "", false); delta {
		t.Errorf("delta with no since")
	}
	if changed, delta, err := i.Delta(&f, f.ContentHash, false); !delta || len(changed) != 0 || err != nil {
		t.Errorf("delta against itself: %v %v", changed, delta)
	}
	for _, s := range [][2]int{{0, 5}, {2001, 5}, {100, 1001}} {
		if err := i.Resize(s[0], s[1]); err == nil || i.t.Cols() != 100 || i.t.Rows() != 5 {
			t.Errorf("Resize(%d, %d): %v, now %dx%d", s[0], s[1], err, i.t.Cols(), i.t.Rows())
		}
	}
	for k := 0; k < 1100; k++ {
		i.Write([]byte("x\r\n"))
	}
	if err := i.Resize(101, 5); err != nil {
		t.Fatal(err)
	}
	if f, _ := i.Read(Request{View: ViewScrollback, Lines: 5000}); f.HistoryLines != 990 || f.Cols != 101 {
		t.Errorf("after resize to 101 cols: %d history lines, %d cols", f.HistoryLines, f.Cols)
	}
}

// TestInstancePanic: a panic in vt during a Write (an injected CSI handler)
// or a Read (a torn terminal) fails that instance only: the call and every
// later one return ErrFailed, the lock is released, and a second instance
// keeps working.
func TestInstancePanic(t *testing.T) {
	bad, _ := New(80, 24, Config{})
	good, _ := New(80, 24, Config{})
	bad.t.InputHandler().RegisterCsiHandler(vt.FunctionID{Final: "z"}, func(*vt.Params) bool { panic("injected") })
	good.t.InputHandler().RegisterCsiHandler(vt.FunctionID{Final: "z"}, func(*vt.Params) bool { return true })
	if n, err := bad.Write([]byte("a\x1b[z")); n != 0 || !errors.Is(err, ErrFailed) {
		t.Fatalf("Write = %d, %v", n, err)
	}
	if err := bad.Err(); !errors.Is(err, ErrFailed) {
		t.Errorf("Err = %v", err)
	}
	if f, err := bad.Read(Request{}); !errors.Is(err, ErrFailed) || f.Lines != nil {
		t.Errorf("Read after failure = %d lines, %v", len(f.Lines), err)
	}
	if _, err := bad.Write([]byte("b")); !errors.Is(err, ErrFailed) {
		t.Errorf("Write after failure: %v", err)
	}
	if err := bad.Resize(10, 10); !errors.Is(err, ErrFailed) {
		t.Errorf("Resize after failure: %v", err)
	}
	if _, err := good.Write([]byte("a\x1b[zb")); err != nil {
		t.Fatal(err)
	}
	if f, err := good.Read(Request{}); err != nil || f.Trimmed()[0] != "ab" || good.Err() != nil {
		t.Errorf("second instance: %q, %v", f.Trimmed()[0], err)
	}

	if err := bad.Err(); strings.Contains(err.Error(), "goroutine") {
		t.Errorf("the sticky error repeats the stack: %v", err)
	}

	for name, call := range map[string]func(*Instance) error{
		"Read":   func(i *Instance) error { _, err := i.Read(Request{}); return err },
		"Resize": func(i *Instance) error { return i.Resize(10, 10) },
		"Delta":  func(i *Instance) error { _, _, err := i.Delta(nil, "", false); return err },
	} {
		torn, _ := New(80, 24, Config{})
		torn.t = nil // a torn terminal; Delta panics on the nil frame
		if err := call(torn); !errors.Is(err, ErrFailed) || !strings.Contains(err.Error(), "goroutine") {
			t.Errorf("%s of a torn instance: %v", name, err)
		}
		if _, err := torn.Write([]byte("x")); !errors.Is(err, ErrFailed) {
			t.Errorf("Write after a failed %s: %v", name, err)
		}
	}
}
