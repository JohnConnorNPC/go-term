package vt

import (
	"fmt"
	"strings"
	"testing"
)

const (
	osc133A = "\x1b]133;A\x07"
	osc133B = "\x1b]133;B\x07"
	osc133C = "\x1b]133;C\x07"
)

// command is what a shell with prompt marks writes for one command line.
func command(cmd, output, end string) string {
	return osc133A + "$ " + osc133B + cmd + "\r\n" + osc133C + output + end
}

func semanticTerminal(cols, scrollback int) *Terminal {
	opts := NewOptions()
	opts.Cols, opts.Rows, opts.Scrollback = cols, 5, scrollback
	return NewTerminal(opts)
}

func wantOutput(t *testing.T, term *Terminal, text string, exit int) {
	t.Helper()
	got, gotExit, ok := term.LastCommandOutput()
	if !ok || got != text || gotExit != exit {
		t.Errorf("LastCommandOutput = %q, %d, %v; want %q, %d", got, gotExit, ok, text, exit)
	}
}

func TestLastCommandOutput(t *testing.T) {
	term := semanticTerminal(20, 100)
	if _, _, ok := term.LastCommandOutput(); ok {
		t.Error("output before any command")
	}
	term.WriteString(command("ls", "a.txt\r\nb.txt\r\n", "\x1b]133;D;0\x07"))
	wantOutput(t, term, "a.txt\nb.txt", 0)

	// A command still running is not the last finished one.
	term.WriteString(command("sleep", "zz", ""))
	wantOutput(t, term, "a.txt\nb.txt", 0)
	term.WriteString("\r\n\x1b]133;D;130\x07")
	wantOutput(t, term, "zz", 130)

	term.WriteString(command("false", "", "\x1b]133;D;1\x07"))
	wantOutput(t, term, "", 1)

	// A line that wrapped comes back as one line; D without a status is -1.
	long := strings.Repeat("x", 25)
	term.WriteString(command("long", long+"\r\nend\r\n", "\x1b]133;D\x07"))
	wantOutput(t, term, long+"\nend", -1)
}

// TestLastCommandOutputWithoutC: a shell that marks only A, B and D has its
// output taken from the line after the command line.
func TestLastCommandOutputWithoutC(t *testing.T) {
	term := semanticTerminal(30, 100)
	term.WriteString(osc133A + "$ " + osc133B + "echo hi\r\nhi\r\n\x1b]133;D;0\x07")
	wantOutput(t, term, "hi", 0)
}

func TestContinuationPromptIsNotACommand(t *testing.T) {
	term := semanticTerminal(30, 100)
	term.WriteString(osc133A + "$ " + osc133B + "for i in 1\r\n\x1b]133;A;k=s\x07> " + osc133B + "do echo $i; done\r\n" + osc133C + "1\r\n\x1b]133;D;0\x07")
	wantOutput(t, term, "1", 0)
	if n := len(term.InputHandler().commands); n != 1 {
		t.Errorf("%d commands, want 1", n)
	}
}

func TestPromptNavigation(t *testing.T) {
	term := semanticTerminal(20, 100)
	var prompts []int
	for i := range 6 {
		b := term.Buffer()
		prompts = append(prompts, b.YBase+b.Y)
		term.WriteString(command("c", strings.Repeat("o\r\n", i+2), "\x1b]133;D;0\x07"))
	}
	b := term.Buffer()
	prompts = append(prompts, b.YBase+b.Y)
	term.WriteString(osc133A + "$ ")
	if b.YDisp != b.YBase {
		t.Fatal("not at the bottom")
	}

	// Up, one prompt at a time, each to the top of the view: every prompt
	// above where the view started, nearest first.
	var want []int
	for i := len(prompts) - 1; i >= 0; i-- {
		if prompts[i] < b.YDisp {
			want = append(want, prompts[i])
		}
	}
	var visited []int
	for term.ScrollToPreviousPrompt() {
		visited = append(visited, b.YDisp)
	}
	if fmt.Sprint(visited) != fmt.Sprint(want) {
		t.Fatalf("visited %v, want %v", visited, want)
	}
	if b.YDisp != prompts[0] {
		t.Errorf("stopped at %d, want the first prompt at %d", b.YDisp, prompts[0])
	}

	// And down again, ending at the bottom.
	for i := 1; term.ScrollToNextPrompt(); i++ {
		if b.YDisp != min(prompts[i], b.YBase) {
			t.Fatalf("next %d: at %d, want %d", i, b.YDisp, min(prompts[i], b.YBase))
		}
	}
	if b.YDisp != b.YBase {
		t.Errorf("did not end at the bottom: %d of %d", b.YDisp, b.YBase)
	}
}

// TestMarksDieWithTheirLines: a command whose lines left the scrollback, or
// were erased, has nothing to give, and its prompt is not jumped to.
func TestMarksDieWithTheirLines(t *testing.T) {
	term := semanticTerminal(20, 10)
	term.WriteString(command("big", strings.Repeat("line\r\n", 30), "\x1b]133;D;0\x07"))
	if _, _, ok := term.LastCommandOutput(); ok {
		t.Error("output whose prompt was trimmed off the scrollback")
	}
	if term.ScrollToPreviousPrompt() {
		t.Error("jumped to a prompt that is gone")
	}

	term.WriteString(command("ls", "x\r\n", "\x1b]133;D;0\x07"))
	wantOutput(t, term, "x", 0)
	term.WriteString("\x1b[2J")
	if _, _, ok := term.LastCommandOutput(); ok {
		t.Error("output survived ED 2")
	}

	term.WriteString(command("ls", "y\r\n", "\x1b]133;D;0\x07"))
	wantOutput(t, term, "y", 0)
	term.WriteString("\x1bc")
	if _, _, ok := term.LastCommandOutput(); ok || len(term.InputHandler().commands) != 0 {
		t.Error("marks survived RIS")
	}
}

// TestMarksOnTheAlternateScreen: a program on the alternate screen marking
// prompts does not lose the shell's, and its own go with the screen.
func TestMarksOnTheAlternateScreen(t *testing.T) {
	term := semanticTerminal(20, 100)
	term.WriteString(command("ls", "x\r\n", "\x1b]133;D;0\x07"))
	term.WriteString("\x1b[?1049h" + command("inner", "y\r\n", "\x1b]133;D;0\x07"))
	wantOutput(t, term, "y", 0)
	term.WriteString("\x1b[?1049l")
	if _, _, ok := term.LastCommandOutput(); ok {
		t.Error("the alternate screen's command outlived it")
	}
	if !term.ScrollToPreviousPrompt() && term.Buffer().YBase > 0 {
		t.Error("lost the shell's prompt")
	}
}
