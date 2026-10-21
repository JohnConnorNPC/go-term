package runner

import "testing"

func scoreAll(t *testing.T, f func(*Probes) Factory, opt Options) (pass, total int) {
	t.Helper()
	pr, err := LoadProbes("../probes.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, ref := range []string{"x53", "xcur"} {
		gs, err := LoadAll("..", ref)
		if err != nil || len(gs) == 0 {
			t.Fatalf("%s: %d goldens, %v", ref, len(gs), err)
		}
		for _, g := range gs {
			r := RunCase(g, f(pr), pr, opt)
			total++
			if r.Status == PASS {
				pass++
			} else if pass > 0 {
				t.Logf("%s %s %s %s", ref, r.ID, r.Status, r.Diff)
			}
		}
	}
	return pass, total
}

func TestStubRed(t *testing.T) {
	stubF := func(*Probes) Factory { f, _ := Lookup("stub"); return f }
	if pass, total := scoreAll(t, stubF, Options{Splits: true}); pass != 0 {
		t.Fatalf("stub passed %d/%d cases, want 0", pass, total)
	} else {
		t.Logf("stub: 0/%d pass (%d/%d fail)", total, total, total)
	}
}

func TestOracleGreen(t *testing.T) {
	if pass, total := scoreAll(t, Oracle, Options{Splits: true}); pass != total {
		t.Fatalf("oracle passed %d/%d cases, want all", pass, total)
	} else {
		t.Logf("oracle: %d/%d pass", pass, total)
	}
}
