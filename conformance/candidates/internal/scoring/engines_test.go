package scoring

import (
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/gitpod"
	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// TestLookup: a candidate comes back behind safe (not the adapter itself), the anchors unwrapped, an unknown name not
// at all; an anchor has no pin.
func TestLookup(t *testing.T) {
	f, ok := Lookup("gitpod", nil)
	if !ok {
		t.Fatal("gitpod not found")
	}
	e, err := f(runner.Config{Cols: 10, Rows: 3, Scrollback: 10, Unicode: "6"})
	if _, raw := e.(*gitpod.Engine); err != nil || raw {
		t.Errorf("gitpod engine %T %v, want it wrapped by safe", e, err)
	}
	for name, want := range map[string]bool{"magnet": true, "xvt": true, "stub": true, "oracle": true, "nope": false} {
		if _, ok := Lookup(name, &runner.Probes{}); ok != want {
			t.Errorf("Lookup(%s) found=%v, want %v", name, ok, want)
		}
	}
	if Pin("stub") != "-" {
		t.Errorf("stub pin %q, want -", Pin("stub"))
	}
}
