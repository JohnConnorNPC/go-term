package vt

import (
	"sync"
	"testing"
)

// TestWcwidthConcurrentFirstUse: terminals on several goroutines may make the
// first Wcwidth call at once; the shared table is built once (run with -race).
func TestWcwidthConcurrentFirstUse(t *testing.T) {
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			term := newTestTerminal(10, 2)
			term.WriteString("中é")
			if w := Wcwidth(0x4e2d); w != 2 {
				t.Errorf("Wcwidth(U+4E2D) = %d", w)
			}
		}()
	}
	wg.Wait()
}
