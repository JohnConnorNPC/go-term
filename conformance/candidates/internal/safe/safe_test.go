package safe

import (
	"errors"
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// fake panics on Write("boom") and otherwise reports its writes in Title.
type fake struct{ got string }

func (f *fake) Write(p []byte) error {
	if string(p) == "boom" {
		panic("boom")
	}
	f.got += string(p)
	return nil
}
func (f *fake) Resize(int, int) error           { return errors.New("no resize") }
func (f *fake) Reset() error                    { return nil }
func (f *fake) Snapshot() (runner.State, error) { return runner.State{P: runner.P{Title: f.got}}, nil }

func engine(t *testing.T) runner.Engine {
	e, err := Wrap(func(runner.Config) (runner.Engine, error) { return &fake{}, nil })(runner.Config{})
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
