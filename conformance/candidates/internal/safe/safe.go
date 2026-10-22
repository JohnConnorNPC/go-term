// Package safe guards a candidate engine (E4): a panic in any call is the error "panic: <v>", a call that runs past
// the deadline is the error "timeout", and an engine that timed out is dead, so every later call fails at once.
// The runner records each error as that case's FAIL (ERROR in the summary) and goes on with the next case.
package safe

import (
	"errors"
	"fmt"
	"time"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// Default is the per-call deadline Wrap uses for d <= 0.
const Default = 10 * time.Second

// Wrap returns f with the factory call and every engine call guarded, each bounded by d.
func Wrap(f runner.Factory, d time.Duration) runner.Factory {
	if d <= 0 {
		d = Default
	}
	return func(c runner.Config) (runner.Engine, error) {
		g := &guard{d: d}
		e, err := run(g, func() (runner.Engine, error) { return f(c) })
		g.e = e
		return g, err
	}
}

type guard struct {
	e    runner.Engine
	d    time.Duration
	dead bool
}

// run calls fn in its own goroutine. A hung fn can't be stopped: it is abandoned and the engine marked dead.
func run[T any](g *guard, fn func() (T, error)) (v T, err error) {
	if g.dead {
		return v, errors.New("dead: an earlier call timed out")
	}
	type res struct {
		v   T
		err error
	}
	done := make(chan res, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				done <- res{err: fmt.Errorf("panic: %v", p)}
			}
		}()
		v, err := fn()
		done <- res{v, err}
	}()
	select {
	case r := <-done:
		return r.v, r.err
	case <-time.After(g.d):
		g.dead = true
		return v, errors.New("timeout")
	}
}

func do(g *guard, fn func() error) error {
	_, err := run(g, func() (bool, error) { return true, fn() })
	return err
}

func (g *guard) Write(p []byte) error            { return do(g, func() error { return g.e.Write(p) }) }
func (g *guard) Resize(c, r int) error           { return do(g, func() error { return g.e.Resize(c, r) }) }
func (g *guard) Reset() error                    { return do(g, g.e.Reset) }
func (g *guard) Snapshot() (runner.State, error) { return run(g, g.e.Snapshot) }
