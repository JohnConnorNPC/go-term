package vt

import (
	"strings"
	"testing"
)

// TestKittyFlags: push, pop, set and query, each screen its own stack.
func TestKittyFlags(t *testing.T) {
	term := newTestTerminal(80, 24)
	var out strings.Builder
	term.OnData = func(s string) { out.WriteString(s) }
	ih := term.InputHandler()
	q := func() string { out.Reset(); term.WriteString("\x1b[?u"); return out.String() }
	if got := q(); got != "\x1b[?0u" {
		t.Errorf("initial %q", got)
	}
	term.WriteString("\x1b[>1u\x1b[>11u")
	if ih.KittyFlags() != 11 {
		t.Errorf("after push %d", ih.KittyFlags())
	}
	term.WriteString("\x1b[=4;2u") // add 4
	if got := q(); got != "\x1b[?15u" {
		t.Errorf("after set %q", got)
	}
	term.WriteString("\x1b[?1049h") // the alternate screen has its own
	if ih.KittyFlags() != 0 {
		t.Errorf("alt screen %d", ih.KittyFlags())
	}
	term.WriteString("\x1b[>8u\x1b[?1049l")
	if ih.KittyFlags() != 15 {
		t.Errorf("back on main %d", ih.KittyFlags())
	}
	term.WriteString("\x1b[<u")
	if ih.KittyFlags() != 1 {
		t.Errorf("after pop %d", ih.KittyFlags())
	}
	term.WriteString("\x1b[<5u")
	if ih.KittyFlags() != 0 {
		t.Errorf("after popping past the bottom %d", ih.KittyFlags())
	}
	term.WriteString("\x1b[>3u")
	ih.ResetKittyKeyboard()
	if ih.KittyFlags() != 0 {
		t.Error("reset")
	}
	// CSI u with no prefix is still restore-cursor.
	term.WriteString("\x1b[5;5H\x1b[s\x1b[1;1H\x1b[u")
	if b := term.Buffer(); b.X != 4 || b.Y != 4 {
		t.Errorf("restore cursor at %d,%d", b.X, b.Y)
	}
}

// TestKittyKey: how keys are reported under each flag.
func TestKittyKey(t *testing.T) {
	k := func(key, code string, keyCode int, mods string) *KeyboardEvent {
		return &KeyboardEvent{Key: key, Code: code, KeyCode: keyCode,
			ShiftKey: strings.Contains(mods, "s"), CtrlKey: strings.Contains(mods, "c"),
			AltKey: strings.Contains(mods, "a"), MetaKey: strings.Contains(mods, "m")}
	}
	for _, tc := range []struct {
		name  string
		ev    *KeyboardEvent
		flags int
		event int
		want  string // "" = legacy / nothing
	}{
		{"off", k("a", "KeyA", 65, ""), 0, KittyPress, ""},
		{"text stays text", k("a", "KeyA", 65, ""), KittyDisambiguate, KittyPress, ""},
		{"shifted text stays text", k("A", "KeyA", 65, "s"), KittyDisambiguate, KittyPress, ""},
		{"escape", k("Escape", "Escape", 27, ""), KittyDisambiguate, KittyPress, "\x1b[27u"},
		{"ctrl+a", k("a", "KeyA", 65, "c"), KittyDisambiguate, KittyPress, "\x1b[97;5u"},
		{"ctrl+shift+a", k("A", "KeyA", 65, "cs"), KittyDisambiguate, KittyPress, "\x1b[97;6u"},
		{"alt+1", k("1", "Digit1", 49, "a"), KittyDisambiguate, KittyPress, "\x1b[49;3u"},
		{"enter stays legacy", k("Enter", "Enter", 13, ""), KittyDisambiguate, KittyPress, ""},
		{"shift+enter", k("Enter", "Enter", 13, "s"), KittyDisambiguate, KittyPress, "\x1b[13;2u"},
		{"ctrl+i is not tab", k("i", "KeyI", 73, "c"), KittyDisambiguate, KittyPress, "\x1b[105;5u"},
		{"arrow stays legacy", k("ArrowUp", "ArrowUp", 38, ""), KittyDisambiguate, KittyPress, ""},
		{"release not asked for", k("a", "KeyA", 65, "c"), KittyDisambiguate, KittyRelease, ""},
		{"release", k("a", "KeyA", 65, "c"), KittyDisambiguate | KittyEventTypes, KittyRelease, "\x1b[97;5:3u"},
		{"arrow release", k("ArrowUp", "ArrowUp", 38, ""), KittyDisambiguate | KittyEventTypes, KittyRelease, "\x1b[1;1:3A"},
		{"repeat", k("Escape", "Escape", 27, ""), KittyDisambiguate | KittyEventTypes, KittyRepeat, "\x1b[27;1:2u"},
		{"all keys: text", k("a", "KeyA", 65, ""), KittyAllKeysAsCodes, KittyPress, "\x1b[97u"},
		{"all keys: enter", k("Enter", "Enter", 13, ""), KittyAllKeysAsCodes, KittyPress, "\x1b[13u"},
		{"all keys: arrow", k("ArrowLeft", "ArrowLeft", 37, ""), KittyAllKeysAsCodes, KittyPress, "\x1b[D"},
		{"all keys: f5", k("F5", "F5", 116, "c"), KittyAllKeysAsCodes, KittyPress, "\x1b[15;5~"},
		{"all keys: shift key", k("Shift", "ShiftLeft", 16, "s"), KittyAllKeysAsCodes, KittyPress, "\x1b[57441;2u"},
		{"alternate key", k("!", "Digit1", 49, "s"), KittyAllKeysAsCodes | KittyAlternateKeys, KittyPress, "\x1b[49:33;2u"},
		{"associated text", k("A", "KeyA", 65, "s"), KittyAllKeysAsCodes | KittyAlternateKeys | KittyAssociatedText, KittyPress, "\x1b[97:65;2;65u"},
	} {
		got, ok := KittyKey(tc.ev, tc.flags, tc.event)
		if !ok {
			got = ""
		}
		if got != tc.want {
			t.Errorf("%s: %q, want %q", tc.name, got, tc.want)
		}
	}
}
