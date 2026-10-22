package norm

import (
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// TestConstants: the exported modes are xterm.js's literal values (the ones equivcheck compares each candidate with).
func TestConstants(t *testing.T) {
	if CMMask != 0x3000000 || CMP16 != 0x1000000 || CMP256 != 0x2000000 || CMRGB != 0x3000000 {
		t.Fatalf("modes %#x %#x %#x %#x", CMMask, CMP16, CMP256, CMRGB)
	}
}

func TestAttr(t *testing.T) {
	colours := []struct {
		mode uint32
		v    int
		want runner.Color
	}{{0, 5, runner.Color{Mode: "d"}}, {0x1000000, 1, runner.Color{Mode: "p16", Value: 1}},
		{0x2000000, 200, runner.Color{Mode: "p256", Value: 200}}, {0x3000000, 0x102030, runner.Color{Mode: "rgb", Value: 0x102030}},
		{0x4000000, 7, runner.Color{Mode: "mode67108864", Value: 7}}}
	for _, c := range colours {
		if got := Attr(Raw{FgMode: c.mode, Fg: c.v}).FG; got != c.want {
			t.Errorf("fg mode %#x: %v, want %v", c.mode, got, c.want)
		}
		if got := Attr(Raw{BgMode: c.mode, Bg: c.v}).BG; got != c.want {
			t.Errorf("bg mode %#x: %v, want %v", c.mode, got, c.want)
		}
	}
	flags := map[string]struct{ raw, want func(*Raw, *runner.Attr) }{
		"b":    {func(r *Raw, _ *runner.Attr) { r.Bold = true }, func(_ *Raw, a *runner.Attr) { a.Bold = true }},
		"i":    {func(r *Raw, _ *runner.Attr) { r.Italic = true }, func(_ *Raw, a *runner.Attr) { a.Italic = true }},
		"d":    {func(r *Raw, _ *runner.Attr) { r.Dim = true }, func(_ *Raw, a *runner.Attr) { a.Dim = true }},
		"u":    {func(r *Raw, _ *runner.Attr) { r.Underline = true }, func(_ *Raw, a *runner.Attr) { a.Underline = true }},
		"bl":   {func(r *Raw, _ *runner.Attr) { r.Blink = true }, func(_ *Raw, a *runner.Attr) { a.Blink = true }},
		"inv":  {func(r *Raw, _ *runner.Attr) { r.Inverse = true }, func(_ *Raw, a *runner.Attr) { a.Inverse = true }},
		"inv0": {func(r *Raw, _ *runner.Attr) { r.Invisible = true }, func(_ *Raw, a *runner.Attr) { a.Invisible = true }},
		"s":    {func(r *Raw, _ *runner.Attr) { r.Strike = true }, func(_ *Raw, a *runner.Attr) { a.Strike = true }},
		"o":    {func(r *Raw, _ *runner.Attr) { r.Overline = true }, func(_ *Raw, a *runner.Attr) { a.Overline = true }},
	}
	for k, f := range flags {
		var r Raw
		want := runner.Attr{FG: runner.Color{Mode: "d"}, BG: runner.Color{Mode: "d"}}
		f.raw(&r, nil)
		f.want(nil, &want)
		if got := Attr(r); got != want {
			t.Errorf("flag %s: %+v, want %+v", k, got, want)
		}
	}
	if !On(0x8000000) || On(0) {
		t.Error("On is not the getter's truthiness")
	}
}
