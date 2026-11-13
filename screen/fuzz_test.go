package screen

import (
	"testing"
	"time"
)

// fuzzOps runs data as Instance operations: a byte op, then its operands.
// op%4 0, 1: write the next n+1 bytes (n the operand byte); 2: resize to
// 1+c%300 x 1+r%100; 3: read (with history and attribute runs), or with
// operand 0xff resize to MaxCols x MaxRows.
func fuzzOps(data []byte) error {
	i, err := New(80, 24, Config{})
	if err != nil {
		return err
	}
	arg := func() int {
		if len(data) == 0 {
			return 0
		}
		b := data[0]
		data = data[1:]
		return int(b)
	}
	for len(data) > 0 {
		switch op := arg(); op % 4 {
		case 0, 1:
			n := min(arg()+1, len(data))
			_, err = i.Write(data[:n])
			data = data[n:]
		case 2:
			err = i.Resize(1+arg()%300, 1+arg()%100)
		default:
			if arg() == 0xff {
				err = i.Resize(MaxCols, MaxRows)
			} else {
				_, err = i.Read(Request{View: ViewScrollback, Lines: 100, Attrs: true})
			}
		}
		if err != nil {
			return err // a panic in vt fails the instance: ErrFailed with the stack
		}
	}
	return nil
}

// FuzzInstance [#7]: random writes, cut anywhere, interleaved with resizes
// (reflow) and reads never panic or hang (10 s per input).
func FuzzInstance(f *testing.F) {
	for _, s := range splitStreams {
		f.Add(append([]byte{0, 0xff}, s...))
		f.Add(append([]byte{2, 3, 2, 0, 7}, s...))
	}
	f.Add([]byte("\x00\x05hello\x03\xff\x00\x03 big\x02\x00\x00\x03\x00"))
	f.Add([]byte("\x00\xff" + "the quick brown fox jumps over the lazy dog\r\n" + "\x02\x01\x00\x02\xc7\x50\x03\x00"))
	f.Fuzz(func(t *testing.T, data []byte) {
		done := make(chan error, 1)
		go func() { done <- fuzzOps(data) }()
		select {
		case err := <-done:
			if err != nil {
				t.Fatal(err)
			}
		case <-time.After(10 * time.Second):
			t.Fatalf("hang: %d bytes of ops still running after 10 s", len(data))
		}
	})
}
