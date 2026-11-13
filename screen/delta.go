package screen

import (
	"container/list"
	"sync"
)

const (
	// FrameCacheSize is how many frames a FrameCache keeps (agent.py
	// _FRAME_CACHE_MAX).
	FrameCacheSize = 16
	// DeltaRatio bounds a delta: at most int(rows*0.6) changed rows
	// (_DELTA_MAX_CHANGED_RATIO). An unchanged frame always qualifies.
	DeltaRatio = 0.6
)

// ChangedRow is one row of a delta: its index in Frame.Lines and its text.
type ChangedRow struct {
	Row  int
	Text string
}

// FrameCache keeps the rows of the last FrameCacheSize frames by content
// hash, least recently remembered evicted first, so a later read can send only
// the rows that changed since a frame the caller holds (Studio's `since`). It
// is safe for concurrent use.
type FrameCache struct {
	mu    sync.Mutex
	order *list.List               // of string hashes, oldest first
	byKey map[string]*list.Element // hash -> element of order
	rows  map[string][]string
}

// NewFrameCache returns an empty cache.
func NewFrameCache() *FrameCache {
	return &FrameCache{order: list.New(), byKey: map[string]*list.Element{}, rows: map[string][]string{}}
}

// Delta compares f with the frame remembered under since and then remembers
// f, as agent.py's screen_text reply does. It returns delta true, with the
// changed rows in order (empty when nothing changed), only when since is a
// remembered hash, the read is not degraded, f.View is ViewScreen, the two
// frames have the same number of rows, and at most int(rows*DeltaRatio) rows
// differ (or none); otherwise delta is false and the caller sends the full
// text. Looking since up does not refresh it; remembering a hash already held
// does. A frame with an empty ContentHash is not remembered.
func (c *FrameCache) Delta(f *Frame, since string, degraded bool) (changed []ChangedRow, delta bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if since != "" && !degraded && f.View == ViewScreen {
		if base, ok := c.rows[since]; ok && len(base) == len(f.Lines) {
			changed = []ChangedRow{}
			for i, l := range f.Lines {
				if l != base[i] {
					changed = append(changed, ChangedRow{Row: i, Text: l})
				}
			}
			// int(len*0.6) without float rounding: len*6/10 floors the same way.
			if len(changed) == 0 || len(changed) <= len(f.Lines)*6/10 {
				delta = true
			} else {
				changed = nil
			}
		}
	}
	c.remember(f.ContentHash, f.Lines)
	return changed, delta
}

func (c *FrameCache) remember(hash string, lines []string) {
	if hash == "" {
		return
	}
	if e, ok := c.byKey[hash]; ok {
		c.order.MoveToBack(e)
		return
	}
	c.byKey[hash] = c.order.PushBack(hash)
	c.rows[hash] = append([]string(nil), lines...)
	for c.order.Len() > FrameCacheSize {
		old := c.order.Remove(c.order.Front()).(string)
		delete(c.byKey, old)
		delete(c.rows, old)
	}
}

// Len is how many frames the cache holds.
func (c *FrameCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.order.Len()
}
