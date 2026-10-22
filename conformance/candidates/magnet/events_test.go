package magnet

import (
	"strings"
	"testing"
)

// TestEvents: one reply per OnData event in order (DA1 then CPR from one write), the titles in order with the last as
// title, and the bell and line-feed counters.
func TestEvents(t *testing.T) {
	st := snap(t, engine(t, 20, 3, "6"), "\x1b[c\x1b[6n")
	if len(st.P.Replies) != 2 || !strings.HasPrefix(st.P.Replies[0], "\x1b[?") || st.P.Replies[1] != "\x1b[1;1R" {
		t.Errorf("replies %q, want DA1 then CPR", st.P.Replies)
	}
	st = snap(t, engine(t, 20, 3, "6"), "\x1b]0;a\a\x1b]2;b\a")
	if strings.Join(st.P.Titles, ",") != "a,b" || st.P.Title != "b" {
		t.Errorf("titles %q title %q, want [a b] b", st.P.Titles, st.P.Title)
	}
	if st := snap(t, engine(t, 20, 3, "6"), "\a\a\n\n"); st.P.Bells != 2 || st.P.LineFeeds != 2 {
		t.Errorf("bells %d lineFeeds %d, want 2 2", st.P.Bells, st.P.LineFeeds)
	}
}
