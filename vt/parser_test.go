package vt

import (
	"fmt"
	"reflect"
	"testing"
)

// feed pushes a string through the UTF-8 decoder and parser like the
// terminal write path does.
func feed(p *Parser, input string) {
	var dec Utf8ToUtf32
	data := make([]uint32, len(input))
	n := dec.Decode([]byte(input), data)
	p.Parse(data, n)
}

func TestParserPrint(t *testing.T) {
	p := NewParser()
	var printed string
	p.SetPrintHandler(func(data []uint32, start, end int) {
		printed += utf32ToString(data, start, end)
	})
	feed(p, "hello world")
	if printed != "hello world" {
		t.Fatalf("printed = %q", printed)
	}
}

func TestParserCSI(t *testing.T) {
	p := NewParser()
	var got []interface{}
	p.RegisterCsiHandler(FunctionID{Final: "m"}, func(params *Params) bool {
		got = params.ToArray()
		return true
	})
	feed(p, "\x1b[1;31;4m")
	want := []interface{}{1, 31, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("params = %v, want %v", got, want)
	}
}

func TestParserCSISubParams(t *testing.T) {
	p := NewParser()
	var got []interface{}
	p.RegisterCsiHandler(FunctionID{Final: "m"}, func(params *Params) bool {
		got = params.ToArray()
		return true
	})
	feed(p, "\x1b[38:2:10:20:30;1m")
	want := []interface{}{38, []int{2, 10, 20, 30}, 1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("params = %v, want %v", got, want)
	}
}

func TestParserCSIPrefixCollect(t *testing.T) {
	p := NewParser()
	hits := 0
	p.RegisterCsiHandler(FunctionID{Prefix: "?", Final: "h"}, func(params *Params) bool {
		hits++
		if params.Params[0] != 25 {
			t.Fatalf("param = %d", params.Params[0])
		}
		return true
	})
	feed(p, "\x1b[?25h")
	if hits != 1 {
		t.Fatalf("hits = %d", hits)
	}
}

func TestParserExecute(t *testing.T) {
	p := NewParser()
	var codes []byte
	p.SetExecuteHandler(0x0a, func() bool { codes = append(codes, 0x0a); return true })
	p.SetExecuteHandler(0x0d, func() bool { codes = append(codes, 0x0d); return true })
	feed(p, "a\r\nb")
	if string(codes) != "\r\n" {
		t.Fatalf("codes = %q", codes)
	}
}

func TestParserOSC(t *testing.T) {
	p := NewParser()
	var title string
	p.RegisterOscHandler(0, NewOscHandler(func(data string) bool {
		title = data
		return true
	}))
	feed(p, "\x1b]0;window title\x07")
	if title != "window title" {
		t.Fatalf("title = %q", title)
	}
	// ST terminator variant
	title = ""
	feed(p, "\x1b]0;other title\x1b\\")
	if title != "other title" {
		t.Fatalf("title = %q", title)
	}
}

func TestParserESC(t *testing.T) {
	p := NewParser()
	hits := 0
	p.RegisterEscHandler(FunctionID{Final: "c"}, func() bool { hits++; return true })
	feed(p, "\x1bc")
	if hits != 1 {
		t.Fatalf("hits = %d", hits)
	}
}

func TestParserDCS(t *testing.T) {
	p := NewParser()
	var payload string
	var params []interface{}
	p.RegisterDcsHandler(FunctionID{Final: "q"}, NewDcsHandler(func(data string, ps *Params) bool {
		payload = data
		params = ps.ToArray()
		return true
	}))
	feed(p, "\x1bP1;2q payload data\x1b\\")
	if payload != " payload data" {
		t.Fatalf("payload = %q", payload)
	}
	if !reflect.DeepEqual(params, []interface{}{1, 2}) {
		t.Fatalf("params = %v", params)
	}
}

func TestParserChunkedSequence(t *testing.T) {
	// a CSI split across arbitrary chunk boundaries must still dispatch once
	seq := "\x1b[38:2:10:20:30;4m"
	for split := 1; split < len(seq); split++ {
		p := NewParser()
		var got []interface{}
		p.RegisterCsiHandler(FunctionID{Final: "m"}, func(params *Params) bool {
			got = params.ToArray()
			return true
		})
		feed(p, seq[:split])
		feed(p, seq[split:])
		want := []interface{}{38, []int{2, 10, 20, 30}, 4}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("split %d: params = %v, want %v", split, got, want)
		}
	}
}

func TestParserPrintAroundSequences(t *testing.T) {
	p := NewParser()
	var printed string
	sgr := 0
	p.SetPrintHandler(func(data []uint32, start, end int) {
		printed += utf32ToString(data, start, end)
	})
	p.RegisterCsiHandler(FunctionID{Final: "m"}, func(params *Params) bool { sgr++; return true })
	feed(p, "red:\x1b[31mtext\x1b[0m.")
	if printed != "red:text." {
		t.Fatalf("printed = %q", printed)
	}
	if sgr != 2 {
		t.Fatalf("sgr = %d", sgr)
	}
}

func TestUtf8DecoderAcrossChunks(t *testing.T) {
	var dec Utf8ToUtf32
	full := "héllo → 世界 🚀"
	raw := []byte(full)
	for split := 1; split < len(raw); split++ {
		dec.Clear()
		var out []rune
		for _, chunk := range [][]byte{raw[:split], raw[split:]} {
			target := make([]uint32, len(chunk)+4)
			n := dec.Decode(chunk, target)
			for i := 0; i < n; i++ {
				out = append(out, rune(target[i])) // #nosec G115 -- ASCII test fixtures
			}
		}
		if string(out) != full {
			t.Fatalf("split %d: got %q want %q", split, string(out), full)
		}
	}
}

func TestParserIdentToString(t *testing.T) {
	p := NewParser()
	ident := p.identifier(FunctionID{Prefix: "?", Final: "h"}, [2]int{0x40, 0x7e})
	if IdentToString(ident) != "?h" {
		t.Fatalf("ident = %q", IdentToString(ident))
	}
}

func TestParserZDM(t *testing.T) {
	// "\x1b[;5H" → params [0, 5] under Zero Default Mode
	p := NewParser()
	var got []interface{}
	p.RegisterCsiHandler(FunctionID{Final: "H"}, func(params *Params) bool {
		got = params.ToArray()
		return true
	})
	feed(p, "\x1b[;5H")
	if !reflect.DeepEqual(got, []interface{}{0, 5}) {
		t.Fatalf("params = %v", got)
	}
}

func ExampleParser() {
	p := NewParser()
	p.SetPrintHandler(func(data []uint32, start, end int) {
		fmt.Print(utf32ToString(data, start, end))
	})
	feed(p, "plain \x1b[31mcolored\x1b[0m text")
	// Output: plain colored text
}

func TestBufferBasics(t *testing.T) {
	opts := NewOptions()
	b := NewBuffer(true, opts, 80, 24)
	b.FillViewportRows(nil)
	if b.Lines.Length() != 24 {
		t.Fatalf("lines = %d", b.Lines.Length())
	}
	// tab stops at multiples of 8
	if b.NextStop(0) != 8 || b.NextStop(8) != 16 || b.PrevStop(9) != 8 {
		t.Fatalf("tab stops wrong: next(0)=%d next(8)=%d prev(9)=%d", b.NextStop(0), b.NextStop(8), b.PrevStop(9))
	}
	// write into a line and translate to string
	line := b.Lines.Get(0)
	attrs := NewAttributeData()
	for i, r := range "hello" {
		line.SetCellFromCodepoint(i, uint32(r), 1, attrs) // #nosec G115 -- ASCII test fixtures
	}
	if s := line.TranslateToString(true, 0, -1); s != "hello" {
		t.Fatalf("line = %q", s)
	}
	// wide char handling
	line.SetCellFromCodepoint(6, 0x4E16, 2, attrs) // 世
	line.SetCellFromCodepoint(7, 0, 0, attrs)      // wide char trailer
	if s := line.TranslateToString(true, 0, -1); s != "hello 世" {
		t.Fatalf("line = %q", s)
	}
	if got := line.GetTrimmedLength(); got != 8 {
		t.Fatalf("trimmed = %d", got)
	}
}

func TestBufferResizeReflowSmaller(t *testing.T) {
	opts := NewOptions()
	opts.Cols = 10
	b := NewBuffer(true, opts, 10, 5)
	b.FillViewportRows(nil)
	b.Y = 3 // keep the cursor off the reflowed line (cursor lines don't reflow)
	attrs := NewAttributeData()
	// a 10-wide line of "abcdefghij" that must wrap when cols -> 5
	line := b.Lines.Get(0)
	for i, r := range "abcdefghij" {
		line.SetCellFromCodepoint(i, uint32(r), 1, attrs) // #nosec G115 -- ASCII test fixtures
	}
	b.Resize(5, 5)
	first := b.Lines.Get(0).TranslateToString(true, 0, -1)
	second := b.Lines.Get(1).TranslateToString(true, 0, -1)
	if first != "abcde" || second != "fghij" {
		t.Fatalf("reflow: %q + %q", first, second)
	}
	if !b.Lines.Get(1).IsWrapped {
		t.Fatal("second line should be wrapped")
	}
	// grow back and verify unwrap
	b.Resize(10, 5)
	joined := b.Lines.Get(0).TranslateToString(true, 0, -1)
	if joined != "abcdefghij" {
		t.Fatalf("unwrap: %q", joined)
	}
}
