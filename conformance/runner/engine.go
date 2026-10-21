package runner

import "sort"

// Config is the terminal an engine is created as (the golden's init).
type Config struct {
	Cols, Rows, Scrollback int
	Unicode                string
}

// Engine is the adapter every scored terminal implements (D7). Snapshot reports the state after the writes so far.
type Engine interface {
	Write(p []byte) error
	Resize(cols, rows int) error
	Reset() error
	Snapshot() (State, error)
}

// Factory creates a fresh engine.
type Factory func(Config) (Engine, error)

// caseAware is unexported on purpose (S5): only engines in this package (the oracle) can be handed the golden they are
// scored on; an engine from another package cannot implement an unexported method, so it never sees a golden.
type caseAware interface{ begin(*Golden) }

var engines = map[string]Factory{"stub": func(Config) (Engine, error) { return stub{}, nil }}

// Register adds a named engine for the CLI.
func Register(name string, f Factory) { engines[name] = f }

// Lookup returns the named engine.
func Lookup(name string) (Factory, bool) { f, ok := engines[name]; return f, ok }

// Engines lists the registered engine names.
func Engines() []string {
	var n []string
	for k := range engines {
		n = append(n, k)
	}
	sort.Strings(n)
	return n
}

// stub does nothing: every snapshot is empty, so every case must FAIL against it.
type stub struct{}

func (stub) Write([]byte) error       { return nil }
func (stub) Resize(int, int) error    { return nil }
func (stub) Reset() error             { return nil }
func (stub) Snapshot() (State, error) { return State{}, nil }
