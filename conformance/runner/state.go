// Package runner scores a terminal engine against the conformance goldens generated from upstream xterm.js headless.
package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Color is a cell colour as the goldens write it: ["d"], ["p16",n], ["p256",n] or ["rgb",0xRRGGBB].
type Color struct {
	Mode  string
	Value int
}

func (c *Color) UnmarshalJSON(b []byte) error {
	var v []any
	if err := json.Unmarshal(b, &v); err != nil || len(v) == 0 || len(v) > 2 {
		return fmt.Errorf("colour %s", b)
	}
	m, ok := v[0].(string)
	if !ok {
		return fmt.Errorf("colour mode %s", b)
	}
	c.Mode, c.Value = m, 0
	if len(v) == 2 {
		f, ok := v[1].(float64)
		if !ok {
			return fmt.Errorf("colour value %s", b)
		}
		c.Value = int(f)
	}
	return nil
}

// Attr is the public (tier P) part of a cell's attributes.
type Attr struct {
	FG, BG                                                                    Color
	Bold, Italic, Dim, Underline, Blink, Inverse, Invisible, Strike, Overline bool
}

func (a *Attr) UnmarshalJSON(b []byte) error {
	var r struct {
		FG                              Color `json:"fg"`
		BG                              Color `json:"bg"`
		B, I, D, U, Bl, Inv, Inv0, S, O int
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&r); err != nil {
		return fmt.Errorf("attr %s: %w", b, err)
	}
	*a = Attr{r.FG, r.BG, r.B == 1, r.I == 1, r.D == 1, r.U == 1, r.Bl == 1, r.Inv == 1, r.Inv0 == 1, r.S == 1, r.O == 1}
	return nil
}

// Cell is one grid cell: Ch is "" for an empty cell and for a wide char's trailing placeholder (W 0).
type Cell struct {
	Ch string
	W  int
	A  Attr
}

// Row is one buffer line.
type Row struct {
	Wrapped bool
	Cells   []Cell
}

// P is the public-API part of a snapshot (xterm.js buffer, modes, options and events).
type P struct {
	Cols         int            `json:"cols"`
	Rows         int            `json:"rows"`
	Buffer       string         `json:"buffer"`
	Cursor       [2]int         `json:"cursor"`
	BaseY        int            `json:"baseY"`
	ViewportY    int            `json:"viewportY"`
	NormalLength int            `json:"normalLength"`
	Modes        map[string]any `json:"modes"`
	Opts         map[string]any `json:"opts"`
	Title        string         `json:"title"`
	Titles       []string       `json:"titles"`
	Bells        int            `json:"bells"`
	Replies      []string       `json:"replies"`
	LineFeeds    int            `json:"lineFeeds"`
}

// State is what an engine reports after a step: P, the normal buffer (scrollback + screen), the alternate screen
// while it is active (nil otherwise), and optional internal values (tier I, reported, never scored). Unobserved lists
// the paths the engine cannot observe (unobservedScalars, modes, opts, modes.<k>, opts.<k>, wrapped, attrs): each one the
// golden has is a FAIL, whatever value the engine reports there.
type State struct {
	P           P
	Normal, Alt []Row
	I           map[string]any
	Unobserved  []string
}
