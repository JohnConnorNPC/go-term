package screen

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestInstanceConcurrent [#7]: one writer (every split stream, resizes) and
// four readers (Read with attribute runs and history, Delta) share an
// Instance; run it with -race. Every frame a reader gets is whole: Rows grid
// rows after its history, the cursor on the grid, the hash of its own text.
func TestInstanceConcurrent(t *testing.T) {
	i, err := New(80, 24, Config{})
	if err != nil {
		t.Fatal(err)
	}
	var stop atomic.Bool
	var reads atomic.Int64
	var wg sync.WaitGroup
	for r := 0; r < 4; r++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			since := ""
			for !stop.Load() {
				f, err := i.Read(Request{View: View(r % 2), Lines: 50, Attrs: r < 2})
				if err != nil {
					t.Error(err)
					return
				}
				if len(f.Lines) != f.Rows+f.HistoryLines || f.Cursor.Row >= f.Rows || f.Cursor.Col >= f.Cols || f.ContentHash != ContentHash(f.Text()) {
					t.Errorf("torn frame: %d lines, %dx%d, %d history, cursor %v", len(f.Lines), f.Cols, f.Rows, f.HistoryLines, *f.Cursor)
					return
				}
				i.Delta(&f, since, false)
				since = f.ContentHash
				reads.Add(1)
			}
		}(r)
	}
	writes := 0
	for deadline := time.Now().Add(time.Second); time.Now().Before(deadline); writes++ {
		for _, s := range splitStreams {
			if _, err := i.Write([]byte(s)); err != nil {
				t.Fatal(err)
			}
		}
		if err := i.Resize(40+writes%60, 10+writes%20); err != nil {
			t.Fatal(err)
		}
	}
	stop.Store(true)
	wg.Wait()
	if reads.Load() == 0 || writes == 0 {
		t.Errorf("%d reads, %d write rounds: nothing ran concurrently", reads.Load(), writes)
	}
	t.Logf("%d write rounds, %d reads", writes, reads.Load())
}
