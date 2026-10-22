package magnet

import (
	vt "github.com/0magnet/xterm-go/vt"
	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// events are the terminal's callbacks, kept from construction on as gen/lib/session.mjs keeps them (a reset op
// clears nothing): one reply per OnData call, in order (no scrolls field: P has none). OnReply stays unset, so magnet
// sends its query replies through OnData too, as xterm.js does through onData.
type events struct {
	replies, titles  []string
	bells, lineFeeds int
}

func (ev *events) listen(t *vt.Terminal) {
	t.OnData = func(s string) { ev.replies = append(ev.replies, s) }
	t.OnTitleChange = func(s string) { ev.titles = append(ev.titles, s) }
	t.OnBell = func() { ev.bells++ }
	t.OnLineFeed = func() { ev.lineFeeds++ }
}

// fill copies the counters into st; the title is the last one set, "" before any (session.mjs titles.at(-1) ?? ”).
func (ev *events) fill(st *runner.State) {
	st.P.Replies, st.P.Titles = append([]string{}, ev.replies...), append([]string{}, ev.titles...)
	st.P.Bells, st.P.LineFeeds = ev.bells, ev.lineFeeds
	if len(ev.titles) > 0 {
		st.P.Title = ev.titles[len(ev.titles)-1]
	}
}
