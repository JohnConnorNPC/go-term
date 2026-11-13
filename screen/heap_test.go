package screen

import (
	"runtime"
	"strings"
	"testing"
)

// heapBound is what 64 MiB of an unterminated OSC, DCS or APC string may add
// to the live heap of an Instance at DefaultPayloadLimit: the string is
// buffered as runes up to the limit (4 MiB), then dropped and its buffer
// released. Without a cap vt would hold 10,000,000 runes (40 MB) per string.
const heapBound = 8 << 20

func liveHeap() uint64 {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapAlloc
}

// TestInstanceHeap [#7]: 64 MiB of each unterminated string, written in 64 KiB
// chunks, keeps the live heap growth under heapBound, sampled every 4 MiB and
// on either side of the limit (just under it the buffer holds ~4.4 MiB).
func TestInstanceHeap(t *testing.T) {
	chunk := []byte(strings.Repeat("x", 64<<10))
	for _, intro := range []string{"\x1b]2;", "\x1bP$q", "\x1b_"} {
		i, err := New(80, 24, Config{})
		if err != nil {
			t.Fatal(err)
		}
		i.t.InputHandler().SetApcHandler(func(string) bool { return true }) // else APC is not buffered at all
		base := liveHeap()
		var peak uint64
		i.Write([]byte(intro))
		for n := 1; n <= 1024; n++ {
			if _, err := i.Write(chunk); err != nil {
				t.Fatal(err)
			}
			if n%64 == 0 || n == 15 || n == 17 { // 15 chunks: just under the limit, 17: just past it
				peak = max(peak, liveHeap()-min(base, liveHeap()))
			}
		}
		t.Logf("%q + 64 MiB: heap growth %.1f MiB (bound %d MiB)", intro, float64(peak)/(1<<20), heapBound>>20)
		if peak >= heapBound {
			t.Errorf("%q: heap grew %d bytes, bound %d", intro, peak, heapBound)
		}
		runtime.KeepAlive(i)
	}
}
