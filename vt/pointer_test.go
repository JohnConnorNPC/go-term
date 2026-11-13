package vt

import (
	"strings"
	"testing"
)

func pointerTerminal(t *testing.T) (*Terminal, *[]string, *strings.Builder) {
	t.Helper()
	term := newTestTerminal(80, 24)
	var shapes []string
	term.OnPointerShape = func(css string) { shapes = append(shapes, css) }
	var reply strings.Builder
	term.OnReply = func(s string) { reply.WriteString(s) }
	return term, &shapes, &reply
}

func last(s []string) string {
	if len(s) == 0 {
		return "<none>"
	}
	return s[len(s)-1]
}

func TestPointerShapeSetPushPop(t *testing.T) {
	term, shapes, _ := pointerTerminal(t)
	for _, c := range []struct{ in, want string }{
		{"\x1b]22;pointer\x1b\\", "pointer"},
		{"\x1b]22;hand2\x07", "pointer"}, // X11 name
		{"\x1b]22;xterm\x07", "text"},    // replaces the top
		{"\x1b]22;>wait,crosshair\x07", "crosshair"},
		{"\x1b]22;<\x07", "wait"},
		{"\x1b]22;<\x07", "text"},
		{"\x1b]22;=fleur\x07", "move"},
		{"\x1b]22;no-such-shape\x07", "move"}, // unknown: ignored
		{"\x1b]22;\x07", ""},                  // back to the default
	} {
		term.WriteString(c.in)
		if got := last(*shapes); got != c.want {
			t.Errorf("%q: %q, want %q", c.in, got, c.want)
		}
	}
}

func TestPointerShapeQuery(t *testing.T) {
	term, _, reply := pointerTerminal(t)
	term.WriteString("\x1b]22;?pointer,bogus,left_ptr,__default__,__current__\x1b\\")
	if got, want := reply.String(), "\x1b]22;1,0,1,text,text\x1b\\"; got != want {
		t.Errorf("query %q, want %q", got, want)
	}
	reply.Reset()
	term.WriteString("\x1b]22;wait\x07\x1b]22;?__current__,__grabbed__\x07")
	if got, want := reply.String(), "\x1b]22;wait,default\x1b\\"; got != want {
		t.Errorf("query %q, want %q", got, want)
	}
}

// TestPointerShapePerScreen: the alternate screen has its own stack, which
// goes with it, and a full reset forgets both.
func TestPointerShapePerScreen(t *testing.T) {
	term, shapes, _ := pointerTerminal(t)
	term.WriteString("\x1b]22;pointer\x07\x1b[?1049h")
	if got := last(*shapes); got != "" {
		t.Errorf("on the alternate screen: %q", got)
	}
	term.WriteString("\x1b]22;grab\x07")
	if got := last(*shapes); got != "grab" {
		t.Errorf("set on the alternate screen: %q", got)
	}
	term.WriteString("\x1b[?1049l")
	if got := last(*shapes); got != "pointer" {
		t.Errorf("back on the main screen: %q", got)
	}
	term.WriteString("\x1b[?1049h")
	if got := last(*shapes); got != "" {
		t.Errorf("a new alternate screen kept %q", got)
	}
	term.WriteString("\x1b[?1049l\x1bc")
	if got := last(*shapes); got != "" {
		t.Errorf("after RIS: %q", got)
	}
}

func TestPointerStackLimit(t *testing.T) {
	term, _, _ := pointerTerminal(t)
	for range 40 {
		term.WriteString("\x1b]22;>wait\x07")
	}
	if n := len(term.InputHandler().pointerMain); n != pointerStackMax {
		t.Errorf("stack grew to %d", n)
	}
}
