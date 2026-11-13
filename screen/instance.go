package screen

import (
	"fmt"
	"sync"

	"github.com/JohnConnorNPC/go-term/vt"
)

const (
	// MaxCols and MaxRows bound an Instance's size: cols 1..MaxCols, rows
	// 1..MaxRows. One column reads as vt.MinimumCols (2).
	MaxCols, MaxRows = 2000, 1000
	// DefaultPayloadLimit is Config.PayloadLimit's default: an OSC, DCS or
	// APC string over 1 MiB code units is dropped (xterm.js takes up to
	// 10,000,000), which bounds what an unterminated one holds.
	DefaultPayloadLimit = 1 << 20
)

// Instance is a terminal safe to use from several goroutines in a server:
// Write, Resize and Read are serialised by one lock, so each Read is a
// consistent snapshot, and Write takes any chunk split (vt keeps a cut
// UTF-8, ESC, CSI, OSC or DCS sequence for the next write).
type Instance struct {
	mu    sync.Mutex
	t     *vt.Terminal
	cfg   Config
	cache *FrameCache
	err   error
}

func checkSize(cols, rows int) error {
	if cols < 1 || cols > MaxCols || rows < 1 || rows > MaxRows {
		return fmt.Errorf("screen: size %dx%d is outside 1..%d x 1..%d", cols, rows, MaxCols, MaxRows)
	}
	return nil
}

// New returns an instance of cols x rows sized and capped by c (NewTerminal).
func New(cols, rows int, c Config) (*Instance, error) {
	if err := checkSize(cols, rows); err != nil {
		return nil, err
	}
	t, err := NewTerminal(cols, rows, c)
	if err != nil {
		return nil, err
	}
	return &Instance{t: t, cfg: c, cache: NewFrameCache()}, nil
}

// do runs f under the lock, unless the instance has failed.
func (i *Instance) do(f func()) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if i.err != nil {
		return i.err
	}
	f()
	return nil
}

// Write feeds pty output to the terminal; it is an io.Writer.
func (i *Instance) Write(p []byte) (int, error) {
	if err := i.do(func() { i.t.Write(p) }); err != nil {
		return 0, err
	}
	return len(p), nil
}

// Resize resizes the terminal as Resize does. A size outside the limits is
// an error and leaves the size unchanged.
func (i *Instance) Resize(cols, rows int) error {
	if err := checkSize(cols, rows); err != nil {
		return err
	}
	return i.do(func() { _ = Resize(i.t, cols, rows, i.cfg) }) // cfg was checked by New
}

// Read is Read of the terminal, under the lock. A failed instance gives a
// zero Frame and its error.
func (i *Instance) Read(r Request) (Frame, error) {
	var f Frame
	err := i.do(func() { f = Read(i.t, r) })
	if err != nil {
		return Frame{}, err
	}
	return f, nil
}

// Delta is the instance's own FrameCache.Delta: f should be a Read of this
// instance.
func (i *Instance) Delta(f *Frame, since string, degraded bool) ([]ChangedRow, bool) {
	return i.cache.Delta(f, since, degraded)
}
