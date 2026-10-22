package xvt

import (
	"io"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
	vt "github.com/charmbracelet/x/vt"
)

// eventsUnobserved: no mode getter and no row citing a mode default, so no start state (GO S-3: the store is the
// unexported `modes ansi.Modes`); no option rows; no line-feed callback; replies only through an io.Pipe that has no
// delivery barrier (E3).
var eventsUnobserved = []string{"lineFeeds", "modes", "opts", "replies"}

// events are the title and bell callbacks, kept from construction on as gen/lib/session.mjs keeps them.
type events struct {
	titles []string
	bells  int
}

// listen drains the reply pipe, so that a reply (DA1, DECRQM, …) never blocks Write, and subscribes the callbacks.
func (ev *events) listen(e *vt.Emulator) {
	go func() { _, _ = io.Copy(io.Discard, e) }()
	e.SetCallbacks(vt.Callbacks{
		Title: func(s string) { ev.titles = append(ev.titles, s) },
		Bell:  func() { ev.bells++ },
	})
}

// fill copies the counters into st and lists the paths x/vt can't show.
func (ev *events) fill(st *runner.State) {
	st.P.Titles, st.P.Bells = append([]string{}, ev.titles...), ev.bells
	if len(ev.titles) > 0 {
		st.P.Title = ev.titles[len(ev.titles)-1]
	}
	st.Unobserved = append(st.Unobserved, eventsUnobserved...)
}
