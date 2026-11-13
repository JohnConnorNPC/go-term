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
