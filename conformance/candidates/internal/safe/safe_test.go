package safe

import (
	"errors"
	"testing"
	"time"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// fake panics on Write("boom"), blocks forever on Write("block"), and otherwise reports its writes in Title.
type fake struct{ got string }

func (f *fake) Write(p []byte) error {
	switch string(p) {
	case "boom":
		panic("boom")
	case "block":
		select {}
	}
	f.got += string(p)
	return nil
}
func (f *fake) Resize(int, int) error           { return errors.New("no resize") }
func (f *fake) Reset() error                    { return nil }
func (f *fake) Snapshot() (runner.State, error) { return runner.State{P: runner.P{Title: f.got}}, nil }

func engine(t *testing.T) runner.Engine {
	e, err := Wrap(func(runner.Config) (runner.Engine, error) { return &fake{}, nil }, 50*time.Millisecond)(runner.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return e
}

func TestPanic(t *testing.T) {
	if err := engine(t).Write([]byte("boom")); err == nil || err.Error() != "panic: boom" {
		t.Fatalf("panic: %v, want panic: boom", err)
	}
}

func TestTimeoutThenDead(t *testing.T) {
	e, t0 := engine(t), time.Now()
	if err := e.Write([]byte("block")); err == nil || err.Error() != "timeout" || time.Since(t0) > time.Second+50*time.Millisecond {
		t.Fatalf("block: %v after %v, want timeout within 50 ms + 1 s", err, time.Since(t0))
	}
	if _, err := e.Snapshot(); err == nil || time.Since(t0) > 2*time.Second {
		t.Fatalf("dead engine snapshot: %v, want a fast error", err)
	}
}

func TestPassthrough(t *testing.T) {
	e := engine(t)
	if err := e.Write([]byte("ab")); err != nil {
		t.Fatal(err)
	}
	st, err := e.Snapshot()
	if err != nil || st.P.Title != "ab" {
		t.Fatalf("passthrough: %q %v, want ab", st.P.Title, err)
	}
	if e.Resize(1, 1) == nil {
		t.Fatal("the engine's Resize error was lost")
	}
}

func TestDefault(t *testing.T) {
	if Default != 10*time.Second {
		t.Fatalf("default deadline %v, want 10s", Default)
	}
}
