package vt

import "testing"

func kbd(ev *KeyboardEvent, appCursor bool) string {
	return EvaluateKeyboardEvent(ev, appCursor, false, false).Key
}

func TestKeyboardBasics(t *testing.T) {
	// enter, tab, backspace, escape
	if got := kbd(&KeyboardEvent{KeyCode: 13, Key: "Enter"}, false); got != "\r" {
		t.Fatalf("enter = %q", got)
	}
	if got := kbd(&KeyboardEvent{KeyCode: 9, Key: "Tab"}, false); got != "\t" {
		t.Fatalf("tab = %q", got)
	}
	if got := kbd(&KeyboardEvent{KeyCode: 9, Key: "Tab", ShiftKey: true}, false); got != "\x1b[Z" {
		t.Fatalf("shift-tab = %q", got)
	}
	if got := kbd(&KeyboardEvent{KeyCode: 8, Key: "Backspace"}, false); got != "\x7f" {
		t.Fatalf("backspace = %q", got)
	}
	if got := kbd(&KeyboardEvent{KeyCode: 27, Key: "Escape"}, false); got != "\x1b" {
		t.Fatalf("escape = %q", got)
	}
}

func TestKeyboardArrows(t *testing.T) {
	if got := kbd(&KeyboardEvent{KeyCode: 38, Key: "ArrowUp"}, false); got != "\x1b[A" {
		t.Fatalf("up = %q", got)
	}
	// application cursor mode
	if got := kbd(&KeyboardEvent{KeyCode: 38, Key: "ArrowUp"}, true); got != "\x1bOA" {
		t.Fatalf("app up = %q", got)
	}
	// modifiers: ctrl+right = \x1b[1;5C
	if got := kbd(&KeyboardEvent{KeyCode: 39, Key: "ArrowRight", CtrlKey: true}, false); got != "\x1b[1;5C" {
		t.Fatalf("ctrl-right = %q", got)
	}
	// shift+alt+left
	if got := kbd(&KeyboardEvent{KeyCode: 37, Key: "ArrowLeft", ShiftKey: true, AltKey: true}, false); got != "\x1b[1;4D" {
		t.Fatalf("shift-alt-left = %q", got)
	}
}

func TestKeyboardCtrlChars(t *testing.T) {
	// ctrl+c = \x03
	if got := kbd(&KeyboardEvent{KeyCode: 67, Key: "c", CtrlKey: true}, false); got != "\x03" {
		t.Fatalf("ctrl-c = %q", got)
	}
	// ctrl+space = NUL
	if got := kbd(&KeyboardEvent{KeyCode: 32, Key: " ", CtrlKey: true}, false); got != "\x00" {
		t.Fatalf("ctrl-space = %q", got)
	}
	// plain char passes through
	if got := kbd(&KeyboardEvent{KeyCode: 65, Key: "a"}, false); got != "a" {
		t.Fatalf("a = %q", got)
	}
	if got := kbd(&KeyboardEvent{KeyCode: 65, Key: "A", ShiftKey: true}, false); got != "A" {
		t.Fatalf("A = %q", got)
	}
}

func TestKeyboardFunctionKeys(t *testing.T) {
	if got := kbd(&KeyboardEvent{KeyCode: 112, Key: "F1"}, false); got != "\x1bOP" {
		t.Fatalf("F1 = %q", got)
	}
	if got := kbd(&KeyboardEvent{KeyCode: 116, Key: "F5"}, false); got != "\x1b[15~" {
		t.Fatalf("F5 = %q", got)
	}
	if got := kbd(&KeyboardEvent{KeyCode: 123, Key: "F12", ShiftKey: true}, false); got != "\x1b[24;2~" {
		t.Fatalf("shift-F12 = %q", got)
	}
	if got := kbd(&KeyboardEvent{KeyCode: 46, Key: "Delete"}, false); got != "\x1b[3~" {
		t.Fatalf("delete = %q", got)
	}
	if got := kbd(&KeyboardEvent{KeyCode: 36, Key: "Home"}, true); got != "\x1bOH" {
		t.Fatalf("home = %q", got)
	}
}

func TestKeyboardAltMeta(t *testing.T) {
	// alt+b -> ESC b (readline word-back)
	if got := kbd(&KeyboardEvent{KeyCode: 66, Key: "b", AltKey: true}, false); got != "\x1bb" {
		t.Fatalf("alt-b = %q", got)
	}
	// alt+enter -> ESC CR
	if got := kbd(&KeyboardEvent{KeyCode: 13, Key: "Enter", AltKey: true}, false); got != "\x1b\r" {
		t.Fatalf("alt-enter = %q", got)
	}
	// shift+pageup scrolls instead of sending
	r := EvaluateKeyboardEvent(&KeyboardEvent{KeyCode: 33, ShiftKey: true}, false, false, false)
	if r.Type != KeyPageUp || r.Key != "" {
		t.Fatalf("shift-pgup = %+v", r)
	}
	// cmd+a on mac selects all
	r = EvaluateKeyboardEvent(&KeyboardEvent{KeyCode: 65, Key: "a", MetaKey: true}, false, true, false)
	if r.Type != KeySelectAll {
		t.Fatalf("cmd-a = %+v", r)
	}
}
