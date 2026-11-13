package screen

import (
	"bytes"
	"testing"
	"time"
)

// fuzzBudget bounds one input's work in grid cells [#11]. Every op costs time
// linear in the grid: at 2000x1000 a Read is ~130 ms CPU (~345 ms under -fuzz
// instrumentation), so without a bound a 42-byte input (a max resize, then 20
// Reads) took 7 s instrumented and could trip the watchdog on a loaded host.
const fuzzBudget = 8_000_000

// fuzzOps runs data as Instance operations: a byte op, then its operands.
// op%4 0, 1: write the next n+1 bytes (n the operand byte); 2: resize to
// 1+c%300 x 1+r%100; 3: read (with history and attribute runs), or with
// operand 0xff resize to MaxCols x MaxRows.
//
// Each op is charged in cells before it runs, and the input ends (nil error)
// at the first op that would take the total past fuzzBudget: a resize costs
// the old plus the new area, a read cols*(rows+100), a write of n bytes
// cols*(n + rows*k) where k counts its ESC (0x1b) and C1 CSI (0x9b) bytes,
// since one escape sequence (ED, DECALN, IL, ...) can touch the whole grid.
// It returns the number of ops run.
func fuzzOps(data []byte) (ops int, err error) {
	i, err := New(80, 24, Config{})
	if err != nil {
		return 0, err
	}
	cols, rows, spent := 80, 24, 0
	charge := func(cells int) bool {
		if spent+cells > fuzzBudget {
			return false
		}
		spent += cells
		return true
	}
	arg := func() int {
		if len(data) == 0 {
			return 0
		}
		b := data[0]
		data = data[1:]
		return int(b)
	}
	for len(data) > 0 {
		switch op := arg(); op % 4 {
		case 0, 1:
			n := min(arg()+1, len(data))
			k := bytes.Count(data[:n], []byte{0x1b}) + bytes.Count(data[:n], []byte{0x9b})
			if !charge(cols * (n + rows*k)) {
				return ops, nil
			}
			_, err = i.Write(data[:n])
			data = data[n:]
		case 2:
			c, r := 1+arg()%300, 1+arg()%100
			if !charge(cols*rows + c*r) {
				return ops, nil
			}
			err, cols, rows = i.Resize(c, r), c, r
		default:
			if arg() == 0xff {
				if !charge(cols*rows + MaxCols*MaxRows) {
					return ops, nil
				}
				err, cols, rows = i.Resize(MaxCols, MaxRows), MaxCols, MaxRows
			} else {
				if !charge(cols * (rows + 100)) {
					return ops, nil
				}
				_, err = i.Read(Request{View: ViewScrollback, Lines: 100, Attrs: true})
			}
		}
		if err != nil {
			return ops, err // a panic in vt fails the instance: ErrFailed with the stack
		}
		ops++
	}
	return ops, nil
}

// FuzzInstance [#7]: random writes, cut anywhere, interleaved with resizes
// (reflow) and reads never panic or hang (60 s per input). A hang never
// returns; a budgeted input takes at most ~0.75 s CPU under -fuzz
// instrumentation (the seed e2ff5de2cef58c25: 0.73 s), and 16 workers on a
// loaded host (load 62) slowed it past the old 10 s [#11], so 60 s leaves
// ~80x for load.
func FuzzInstance(f *testing.F) {
	for _, s := range splitStreams {
		f.Add(append([]byte{0, 0xff}, s...))
		f.Add(append([]byte{2, 3, 2, 0, 7}, s...))
	}
	f.Add([]byte("\x00\x05hello\x03\xff\x00\x03 big\x02\x00\x00\x03\x00"))
	f.Add([]byte("\x00\xff" + "the quick brown fox jumps over the lazy dog\r\n" + "\x02\x01\x00\x02\xc7\x50\x03\x00"))
	f.Fuzz(func(t *testing.T, data []byte) {
		done := make(chan error, 1)
		go func() { _, err := fuzzOps(data); done <- err }()
		select {
		case err := <-done:
			if err != nil {
				t.Fatal(err)
			}
		case <-time.After(60 * time.Second):
			t.Fatalf("hang: %d bytes of ops still running after 60 s", len(data))
		}
	})
}

// TestFuzzBudget [#11]: fuzzOps stops an input at the first op past
// fuzzBudget and runs every op of one under it, counted by ops run.
func TestFuzzBudget(t *testing.T) {
	reads20 := []byte("\x03\xff") // max resize, then 20 Reads at 2,200,000 cells each
	for range 20 {
		reads20 = append(reads20, 3, 0)
	}
	decaln := append([]byte("\x03\xff\x00\xff"), bytes.Repeat([]byte("\x1b#8"), 85)...) // one write, 85 ESCs
	small := bytes.Repeat([]byte("\x00\x00x"), 200)                                     // 200 one-byte writes at 80x24
	for _, c := range []struct {
		name string
		data []byte
		ops  int
	}{
		{"seed e2ff5de2cef58c25 (6.4M cells, under)", []byte("\xb5\x00\x10\x00\x00\x7f\xff\xff\xff_\x1b\\"), 5},
		{"max resize + 20 Reads (over at the 3rd Read)", reads20, 3},
		{"max resize + 85 DECALN in one write (over)", decaln, 1},
		{"200 small writes (under)", small, 200},
	} {
		ops, err := fuzzOps(c.data)
		if err != nil || ops != c.ops {
			t.Errorf("%s: %d ops run, err %v; want %d", c.name, ops, err, c.ops)
		}
	}
}
