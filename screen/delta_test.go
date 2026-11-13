package screen

import (
	"fmt"
	"sync"
	"testing"
)

// frame builds a screen frame of n rows "r0".."r<n-1>", the first k of them
// replaced by "x<i>".
func frame(n, k int) *Frame {
	f := &Frame{View: ViewScreen, Rows: n}
	for i := 0; i < n; i++ {
		l := fmt.Sprintf("r%d", i)
		if i < k {
			l = fmt.Sprintf("x%d", i)
		}
		f.Lines = append(f.Lines, l)
	}
	f.ContentHash = ContentHash(f.Text())
	return f
}

// TestFrameCacheBoundary pins the 60 % rule both ways: int(rows*0.6) changed
// rows give a delta, one more gives the full grid.
func TestFrameCacheBoundary(t *testing.T) {
	for _, tc := range []struct {
		rows, changed int
		want          bool
	}{{10, 6, true}, {10, 7, false}, {24, 14, true}, {24, 15, false}, {1, 0, true}, {1, 1, false}, {5, 3, true}, {5, 4, false}} {
		want := tc.want
		if py := tc.changed == 0 || tc.changed <= int(float64(tc.rows)*DeltaRatio); py != want {
			t.Errorf("rows %d changed %d: table says %v, Python's int(rows*0.6) rule %v", tc.rows, tc.changed, want, py)
		}
		c := NewFrameCache()
		base := frame(tc.rows, 0)
		c.Delta(base, "", false)
		rows, delta := c.Delta(frame(tc.rows, tc.changed), base.ContentHash, false)
		if delta != want {
			t.Errorf("rows %d changed %d: delta %v, want %v", tc.rows, tc.changed, delta, want)
		}
		if delta && len(rows) != tc.changed {
			t.Errorf("rows %d changed %d: %d changed rows", tc.rows, tc.changed, len(rows))
		}
		if !delta && rows != nil {
			t.Errorf("rows %d changed %d: full grid with changed rows %v", tc.rows, tc.changed, rows)
		}
	}
}

func TestFrameCacheDelta(t *testing.T) {
	c := NewFrameCache()
	base := frame(4, 0)
	if rows, delta := c.Delta(base, "", false); delta || rows != nil {
		t.Errorf("no since: delta %v", delta)
	}
	next := frame(4, 1)
	rows, delta := c.Delta(next, base.ContentHash, false)
	if !delta || len(rows) != 1 || rows[0] != (ChangedRow{0, "x0"}) {
		t.Errorf("small change: delta %v rows %v", delta, rows)
	}
	if rows, delta := c.Delta(next, next.ContentHash, false); !delta || rows == nil || len(rows) != 0 {
		t.Errorf("unchanged: delta %v rows %#v, want an empty delta", delta, rows)
	}
	for name, run := range map[string]func() bool{
		"unknown since": func() bool { _, d := c.Delta(next, ContentHash("never"), false); return d },
		"degraded":      func() bool { _, d := c.Delta(next, base.ContentHash, true); return d },
		"scrollback": func() bool {
			f := *next
			f.View = ViewScrollback
			_, d := c.Delta(&f, base.ContentHash, false)
			return d
		},
		"row count": func() bool { _, d := c.Delta(frame(5, 0), base.ContentHash, false); return d },
	} {
		if run() {
			t.Errorf("%s: delta, want the full grid", name)
		}
	}
}

// TestFrameCacheLRU: 16 frames are kept; the 17th evicts the least recently
// remembered; remembering a held hash refreshes it, looking it up does not.
func TestFrameCacheLRU(t *testing.T) {
	c := NewFrameCache()
	fs := make([]*Frame, 18)
	for i := range fs {
		fs[i] = frame(3, 0)
		fs[i].Lines[2] = fmt.Sprint("frame", i)
		fs[i].ContentHash = ContentHash(fs[i].Text())
	}
	for i := 0; i < 16; i++ {
		c.Delta(fs[i], "", false)
	}
	c.Delta(fs[0], "", false)                 // refresh frame 0: frame 1 is now the oldest
	c.Delta(fs[16], fs[2].ContentHash, false) // a lookup of frame 2 does not refresh it
	if c.Len() != 16 {
		t.Fatalf("len %d", c.Len())
	}
	known := func(i int) bool {
		_, d := c.Delta(&Frame{View: ViewScreen, Lines: fs[i].Lines}, fs[i].ContentHash, false)
		return d
	}
	if known(1) {
		t.Error("frame 1 not evicted by the 17th")
	}
	if !known(0) || !known(2) || !known(16) {
		t.Error("frames 0, 2 or 16 lost")
	}
	c.Delta(fs[17], "", false) // evicts frame 2, the oldest remembered now
	if known(2) {
		t.Error("frame 2 was refreshed by a lookup")
	}
	c.Delta(&Frame{View: ViewScreen, Lines: []string{"a"}}, "", false)
	if c.Len() != 16 {
		t.Errorf("an empty hash was remembered: len %d", c.Len())
	}
}

func TestFrameCacheConcurrent(t *testing.T) {
	c := NewFrameCache()
	var wg sync.WaitGroup
	for g := 0; g < 4; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				f := frame(10, (g+i)%10)
				c.Delta(f, ContentHash(fmt.Sprint(i)), false)
			}
		}(g)
	}
	wg.Wait()
	if c.Len() > FrameCacheSize {
		t.Errorf("len %d", c.Len())
	}
}
