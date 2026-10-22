// Package safe guards a candidate engine (E4): a panic in any call is the error "panic: <v>". The runner records
// the error as that case's FAIL (ERROR in the summary) and goes on with the next case.
package safe

import (
	"fmt"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// Wrap returns f with the factory call and every engine call guarded.
func Wrap(f runner.Factory) runner.Factory {
	return func(c runner.Config) (runner.Engine, error) {
		g := &guard{}
		e, err := run(g, func() (runner.Engine, error) { return f(c) })
		g.e = e
		return g, err
	}
}

type guard struct {
	e runner.Engine
}

// run calls fn and turns a panic into an error.
func run[T any](g *guard, fn func() (T, error)) (v T, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("panic: %v", p)
		}
	}()
	return fn()
}

func do(g *guard, fn func() error) error {
	_, err := run(g, func() (bool, error) { return true, fn() })
	return err
}

func (g *guard) Write(p []byte) error            { return do(g, func() error { return g.e.Write(p) }) }
func (g *guard) Resize(c, r int) error           { return do(g, func() error { return g.e.Resize(c, r) }) }
func (g *guard) Reset() error                    { return do(g, g.e.Reset) }
func (g *guard) Snapshot() (runner.State, error) { return run(g, g.e.Snapshot) }
