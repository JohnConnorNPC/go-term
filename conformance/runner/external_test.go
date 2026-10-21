package runner_test

import (
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// outsider is an engine from another package that asks for the golden the way the oracle does.
type outsider struct {
	runner.Engine
	got *bool
}

// Both spellings: if the runner's hook were exported again (Begin), this engine would be handed the golden.
func (o outsider) begin(*runner.Golden) { *o.got = true }
func (o outsider) Begin(*runner.Golden) { *o.got = true }

// TestGoldenHookPrivate (S5): an engine outside package runner is never handed the golden it is scored on.
func TestGoldenHookPrivate(t *testing.T) {
	pr, err := runner.LoadProbes("../probes.json")
	if err != nil {
		t.Fatal(err)
	}
	g, err := runner.Load("../golden/x53/c0/smoke-text.json")
	if err != nil {
		t.Fatal(err)
	}
	stub, _ := runner.Lookup("stub")
	called := false
	f := func(c runner.Config) (runner.Engine, error) { e, _ := stub(c); return outsider{e, &called}, nil }
	runner.RunCase(g, f, pr, runner.Options{})
	if called {
		t.Fatal("an external engine was handed the golden")
	}
}
