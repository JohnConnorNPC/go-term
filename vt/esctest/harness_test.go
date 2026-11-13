// Package esctest runs conformance cases adapted from go-te
// (https://github.com/rcarmo/go-te, MIT License, Copyright (c) 2026 Rui
// Carmo), whose tests are themselves ported from esctest2
// (https://github.com/ThomasDickey/esctest2) by George Nachman and Thomas E.
// Dickey. Each case keeps the "From esctest2/..." line naming its origin.
//
// This file is the harness: the helpers those cases call, implemented over
// vt.Terminal. The cases are in cases_test.go; the ones left out, and why,
// are listed at its top.
package esctest

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/vt"
)

const (
	ControlESC = "\x1b"
	ControlCSI = "\x1b["
	ControlOSC = "\x1b]"
	ControlST  = "\x1b\\"
	ControlLF  = "\n"
	ControlCR  = "\r"
	ControlBS  = "\b"
	ControlHT  = "\t"
	ControlFF  = "\f"
	ControlVT  = "\v"

	EscDL  = "M"
	EscIL  = "L"
	EscNEL = "E"
)

// Screen wraps a terminal and captures what it sends back.
type Screen struct {
	term              *vt.Terminal
	WriteProcessInput func(string)
	Title             string
	IconName          string
}

// Stream feeds a Screen.
type Stream struct{ screen *Screen }

func NewScreen(cols, rows int) *Screen {
	opts := vt.NewOptions()
	opts.Cols, opts.Rows = cols, rows
	opts.WindowOptions.SetWinLines = true // xterm.js gates DECCOLM behind this
	s := &Screen{term: vt.NewTerminal(opts)}
	s.term.OnData = func(data string) {
		if s.WriteProcessInput != nil {
			s.WriteProcessInput(data)
		}
	}
	s.term.OnTitleChange = func(title string) { s.Title = title }
	return s
}

func NewStream(screen *Screen, _ bool) *Stream { return &Stream{screen: screen} }

func (s *Stream) Feed(data string) error {
	s.screen.term.WriteString(data)
	return nil
}

type esctestPoint struct{ X, Y int }

type esctestSize struct{ Width, Height int }

type esctestRect struct{ Left, Top, Right, Bottom int }

const (
	esctestModeDECAWM            = 7
	esctestModeDECOM             = 6
	esctestModeIRM               = 4
	esctestModeLNM               = 20
	esctestModeAllow80To132      = 40
	esctestModeAltBuf            = 47
	esctestModeOptAltBuf         = 1047
	esctestModeOptAltBufCursor   = 1049
	esctestModeSaveRestoreCursor = 1048
	esctestModeReverseWrapInline = 45
	esctestModeReverseWrapExtend = 1045
	esctestModeDECCKM            = 1
	esctestModeDECNKM            = 66
	esctestModeDECBKM            = 67
	esctestModeDECNCSM           = 95
	esctestXtermReverseWrap      = 383
	esctestModeDECRLM            = 34
	esctestModeMoreFix           = 41
	esctestModeDECSCLM           = 4
	esctestModeDECARM            = 8
	esctestModeDECPFF            = 18
	esctestModeDECPEX            = 19
	esctestModeDECNRCM           = 42
	esctestModeDECKBUM           = 68
	esctestModeDECHCCM           = 60
	esctestModeDECAAM            = 100
	esctestModeDECCANSM          = 101
	esctestModeDECNULM           = 102
	esctestModeDECHDPXM          = 103
	esctestModeDECESKM           = 104
	esctestModeDECOSCNM          = 106

	esctestTitleSetHex    = 0
	esctestTitleQueryHex  = 1
	esctestTitleSetUTF8   = 2
	esctestTitleQueryUTF8 = 3
)

func esctestReverseWraparoundMode() int {
	if esctestXtermReverseWrap >= 383 {
		return esctestModeReverseWrapExtend
	}
	return esctestModeReverseWrapInline
}

func esctestWrite(t *testing.T, stream *Stream, data string) {
	t.Helper()
	if err := stream.Feed(data); err != nil {
		t.Fatal(err)
	}
}

func esctestJoinParams(params ...int) string {
	parts := make([]string, len(params))
	for i, p := range params {
		parts[i] = strconv.Itoa(p)
	}
	return strings.Join(parts, ";")
}

// csi writes CSI <prefix><params><final>.
func csi(t *testing.T, stream *Stream, prefix, final string, params ...int) {
	t.Helper()
	esctestWrite(t, stream, ControlCSI+prefix+esctestJoinParams(params...)+final)
}

func esctestCUP(t *testing.T, stream *Stream, p esctestPoint) { csi(t, stream, "", "H", p.Y, p.X) }

func optional(row, col *int) string {
	var parts []string
	for _, v := range []*int{row, col} {
		if v == nil {
			parts = append(parts, "")
		} else {
			parts = append(parts, strconv.Itoa(*v))
		}
	}
	for len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return strings.Join(parts, ";")
}

func esctestCUPParams(t *testing.T, stream *Stream, row, col *int) {
	esctestWrite(t, stream, ControlCSI+optional(row, col)+"H")
}

func esctestHVPParams(t *testing.T, stream *Stream, row, col *int) {
	esctestWrite(t, stream, ControlCSI+optional(row, col)+"f")
}

func esctestCUU(t *testing.T, stream *Stream, params ...int) { csi(t, stream, "", "A", params...) }
func esctestCUD(t *testing.T, stream *Stream, params ...int) { csi(t, stream, "", "B", params...) }
func esctestCUF(t *testing.T, stream *Stream, params ...int) { csi(t, stream, "", "C", params...) }
func esctestCUB(t *testing.T, stream *Stream, params ...int) { csi(t, stream, "", "D", params...) }
func esctestCNL(t *testing.T, stream *Stream, params ...int) { csi(t, stream, "", "E", params...) }
func esctestCPL(t *testing.T, stream *Stream, params ...int) { csi(t, stream, "", "F", params...) }
func esctestCHA(t *testing.T, stream *Stream, params ...int) { csi(t, stream, "", "G", params...) }
func esctestCHT(t *testing.T, stream *Stream, params ...int) { csi(t, stream, "", "I", params...) }
func esctestED(t *testing.T, stream *Stream, params ...int)  { csi(t, stream, "", "J", params...) }
func esctestEL(t *testing.T, stream *Stream, params ...int)  { csi(t, stream, "", "K", params...) }
func esctestDCH(t *testing.T, stream *Stream, params ...int) { csi(t, stream, "", "P", params...) }
func esctestICH(t *testing.T, stream *Stream, params ...int) { csi(t, stream, "", "@", params...) }
func esctestSU(t *testing.T, stream *Stream, params ...int)  { csi(t, stream, "", "S", params...) }
func esctestSD(t *testing.T, stream *Stream, params ...int)  { csi(t, stream, "", "T", params...) }
func esctestECH(t *testing.T, stream *Stream, params ...int) { csi(t, stream, "", "X", params...) }
func esctestCBT(t *testing.T, stream *Stream, params ...int) { csi(t, stream, "", "Z", params...) }
func esctestREP(t *testing.T, stream *Stream, params ...int) { csi(t, stream, "", "b", params...) }
func esctestVPA(t *testing.T, stream *Stream, params ...int) { csi(t, stream, "", "d", params...) }
func esctestVPR(t *testing.T, stream *Stream, params ...int) { csi(t, stream, "", "e", params...) }
func esctestHPA(t *testing.T, stream *Stream, params ...int) { csi(t, stream, "", "`", params...) }
func esctestHPR(t *testing.T, stream *Stream, params ...int) { csi(t, stream, "", "a", params...) }
func esctestTBC(t *testing.T, stream *Stream, params ...int) { csi(t, stream, "", "g", params...) }
func esctestSM(t *testing.T, stream *Stream, params ...int)  { csi(t, stream, "", "h", params...) }
func esctestRM(t *testing.T, stream *Stream, params ...int)  { csi(t, stream, "", "l", params...) }
func esctestDECSTBM(t *testing.T, stream *Stream, params ...int) {
	csi(t, stream, "", "r", params...)
}
func esctestDECSET(t *testing.T, stream *Stream, modes ...int)   { csi(t, stream, "?", "h", modes...) }
func esctestDECRESET(t *testing.T, stream *Stream, modes ...int) { csi(t, stream, "?", "l", modes...) }
func esctestDECSED(t *testing.T, stream *Stream, params ...int)  { csi(t, stream, "?", "J", params...) }
func esctestDECSEL(t *testing.T, stream *Stream, params ...int)  { csi(t, stream, "?", "K", params...) }
func esctestDECSCA(t *testing.T, stream *Stream, params ...int)  { csi(t, stream, "", "\"q", params...) }
func esctestDECSTR(t *testing.T, stream *Stream)                 { csi(t, stream, "", "!p") }

func esctestDECRQM(t *testing.T, stream *Stream, mode int, dec bool) {
	prefix := ""
	if dec {
		prefix = "?"
	}
	csi(t, stream, prefix, "$p", mode)
}

func esctestIND(t *testing.T, stream *Stream)    { esctestWrite(t, stream, ControlESC+"D") }
func esctestNEL(t *testing.T, stream *Stream)    { esctestWrite(t, stream, ControlESC+"E") }
func esctestHTS(t *testing.T, stream *Stream)    { esctestWrite(t, stream, ControlESC+"H") }
func esctestRI(t *testing.T, stream *Stream)     { esctestWrite(t, stream, ControlESC+"M") }
func esctestRIS(t *testing.T, stream *Stream)    { esctestWrite(t, stream, ControlESC+"c") }
func esctestDECSC(t *testing.T, stream *Stream)  { esctestWrite(t, stream, ControlESC+"7") }
func esctestDECRC(t *testing.T, stream *Stream)  { esctestWrite(t, stream, ControlESC+"8") }
func esctestDECALN(t *testing.T, stream *Stream) { esctestWrite(t, stream, ControlESC+"#8") }

func esctestChangeWindowTitle(t *testing.T, stream *Stream, title string) {
	esctestWrite(t, stream, ControlOSC+"2;"+title+ControlST)
}

func esctestChangeIconTitle(t *testing.T, stream *Stream, title string) {
	esctestWrite(t, stream, ControlOSC+"1;"+title+ControlST)
}

func esctestGetWindowTitle(screen *Screen) string { return screen.Title }
func esctestGetIconTitle(screen *Screen) string   { return screen.IconName }

func esctestCaptureResponse(screen *Screen, fn func()) string {
	var response string
	prev := screen.WriteProcessInput
	screen.WriteProcessInput = func(data string) { response += data }
	fn()
	screen.WriteProcessInput = prev
	return response
}

func esctestReadCSI(t *testing.T, response, final string, prefix rune) []int {
	t.Helper()
	payload, ok := strings.CutPrefix(response, ControlCSI)
	if !ok {
		t.Fatalf("expected a CSI response, got %q", response)
	}
	if prefix != 0 {
		if payload == "" || rune(payload[0]) != prefix {
			t.Fatalf("expected CSI prefix %q, got %q", prefix, response)
		}
		payload = payload[1:]
	}
	payload, ok = strings.CutSuffix(payload, final)
	if !ok {
		t.Fatalf("expected CSI final %q, got %q", final, response)
	}
	if payload == "" {
		return nil
	}
	var params []int
	for _, part := range strings.Split(payload, ";") {
		v, err := strconv.Atoi(part)
		if err != nil {
			t.Fatalf("invalid CSI param %q in %q", part, response)
		}
		params = append(params, v)
	}
	return params
}

// esctestGetCursorPosition asks with DSR 6, as esctest does, so origin mode
// is reflected the way an application sees it.
func esctestGetCursorPosition(screen *Screen) esctestPoint {
	r := esctestCaptureResponse(screen, func() { screen.term.WriteString(ControlCSI + "6n") })
	r = strings.TrimSuffix(strings.TrimPrefix(r, ControlCSI), "R")
	y, x, _ := strings.Cut(r, ";")
	yi, errY := strconv.Atoi(y)
	xi, errX := strconv.Atoi(x)
	if errY != nil || errX != nil {
		return esctestPoint{}
	}
	return esctestPoint{X: xi, Y: yi}
}

func esctestGetScreenSize(screen *Screen) esctestSize {
	return esctestSize{Width: screen.term.Cols(), Height: screen.term.Rows()}
}

func esctestAssertEQ(t *testing.T, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func esctestEmpty() string { return " " }
func esctestBlank() string { return " " }

func esctestAssertScreenCharsInRectEqual(t *testing.T, screen *Screen, rect esctestRect, expected []string) {
	t.Helper()
	if rows := rect.Bottom - rect.Top + 1; rows != len(expected) {
		t.Fatalf("expected %d rows, got %d", rows, len(expected))
	}
	b := screen.term.Buffer()
	cell := &vt.CellData{}
	for row := rect.Top; row <= rect.Bottom; row++ {
		line := b.Lines.Get(b.YBase + row - 1)
		var sb strings.Builder
		for col := rect.Left; col <= rect.Right; col++ {
			ch := ""
			if line != nil && col-1 < line.Length {
				ch = line.LoadCell(col-1, cell).GetChars()
			}
			if ch == "" {
				ch = " "
			}
			sb.WriteString(ch)
		}
		if got, want := sb.String(), expected[row-rect.Top]; got != want {
			t.Fatalf("row %d: expected %q, got %q", row, want, got)
		}
	}
}

func esctestSMTitle(t *testing.T, stream *Stream, params ...int) { csi(t, stream, ">", "t", params...) }
func esctestRMTitle(t *testing.T, stream *Stream, params ...int) { csi(t, stream, ">", "T", params...) }

// Private mode numbers the way pyte stores them, shifted left by five; the
// cases shift them back.
const (
	ModeDECCOLM = 3 << 5
	ModeDECSCNM = 5 << 5
	ModeDECTCEM = 25 << 5
)
