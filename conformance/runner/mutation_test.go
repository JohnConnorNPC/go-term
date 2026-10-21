package runner

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// pinned is the oracle bound to the pristine golden: begin is shadowed, so it never sees the mutated expected side.
type pinned struct{ *oracle }

func (pinned) begin(*Golden) {}

var (
	freeKeys = regexp.MustCompile(`^(steps\.\*\.P\.(modes|opts)|B\.t1)\.[^.]+$`)
	// Not scored on purpose: identity/pins (STALE has its own tests), display text, tier I, split bookkeeping.
	unscored = regexp.MustCompile(`^(id|ref|pins|inSha|init|split)(\.|$)|^steps\.\*\.(op|text|I)(\.|$)|^steps\.\*\.attrs\.\*\.I(\.|$)`)
)

func decode(t *testing.T, b []byte) map[string]any {
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	var m map[string]any
	if err := d.Decode(&m); err != nil {
		t.Fatal(err)
	}
	return m
}

// TestMutationExactlyOne (M1): for every (category, leaf kind), the first mutable occurrence in sorted order is changed
// on the expected side only; the pinned oracle still replays the pristine golden. Exactly that case must fail and every
// other case must pass (their results are the pristine baseline: a mutation touches one golden only).
func TestMutationExactlyOne(t *testing.T) {
	pr, err := LoadProbes("../probes.json")
	gs, err2 := LoadAll("..", "x53")
	if err != nil || err2 != nil || len(gs) == 0 {
		t.Fatalf("probes %v, goldens %d %v", err, len(gs), err2)
	}
	type occ struct {
		i    int
		path string
	}
	raw, base, cand := make([][]byte, len(gs)), make([]Status, len(gs)), map[string][]occ{}
	for i, g := range gs {
		if raw[i], err = os.ReadFile(filepath.Join("..", "golden", "x53", g.ID+".json")); err != nil {
			t.Fatal(err)
		}
		if base[i] = RunCase(g, Oracle(pr), pr, Options{}).Status; base[i] != PASS { // control arm
			t.Fatalf("pristine %s: %s", g.ID, base[i])
		}
		walk(decode(t, raw[i]), "", func(p string) {
			if k := leafKind(p); !unscored.MatchString(k) {
				cand[g.Category()+" "+k] = append(cand[g.Category()+" "+k], occ{i, p})
			}
		})
	}
	keys := make([]string, 0, len(cand))
	for k := range cand {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	rejected, perCat, kinds, unmutable := 0, map[string]int{}, map[string]bool{}, []string{}
	for _, key := range keys {
		cat, kind, done := strings.SplitN(key, " ", 2)[0], strings.SplitN(key, " ", 2)[1], false
		for _, o := range cand[key] {
			m := decode(t, raw[o.i])
			if !mutateAt(m, o.path) {
				continue
			}
			mb, _ := json.Marshal(m)
			want, err := Parse(mb)
			done = true
			if err != nil {
				rejected++
				t.Errorf("%s: mutating %s in %s is rejected by the parser: %v", key, o.path, gs[o.i].ID, err)
				break
			}
			pristine := gs[o.i]
			f := func(Config) (Engine, error) { return pinned{&oracle{pr: pr, step: -1, g: pristine}}, nil }
			res := append([]Status(nil), base...)
			res[o.i] = RunCase(want, f, pr, Options{}).Status
			fails := 0
			for _, s := range res {
				if s != PASS {
					fails++
				}
			}
			if fails != 1 || res[o.i] == PASS {
				t.Errorf("%s (%s in %s): %d cases failed, mutated case %s; want exactly that one", key, o.path, pristine.ID, fails, res[o.i])
			}
			perCat[cat]++
			kinds[kind] = true
			break
		}
		if !done {
			unmutable = append(unmutable, key)
		}
	}
	if rejected > 0 {
		t.Errorf("%d mutations rejected by the parser, want 0", rejected)
	}
	var cats []string
	for c := range RequiredCategories {
		cats = append(cats, c)
	}
	sort.Strings(cats)
	for _, c := range cats {
		t.Logf("coverage %s kinds=%d", c, perCat[c])
	}
	if len(perCat) != len(RequiredCategories) {
		t.Errorf("mutated categories %v, want %v", perCat, cats)
	}
	for _, k := range requiredKinds {
		if !kinds[k] {
			t.Errorf("required leaf kind %s was mutated in no category", k)
		}
	}
	t.Logf("pairs=%d kinds=%d rejected=%d unmutable=%d %v", len(keys)-len(unmutable), len(kinds), rejected, len(unmutable), unmutable)
}

func TestNoSkips(t *testing.T) {
	pat := "t." + "Skip"
	err := filepath.Walk("..", func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(p, ".go") || filepath.Base(p) == "mutation_test.go" {
			return err
		}
		if b, _ := os.ReadFile(p); bytes.Contains(b, []byte(pat)) {
			t.Errorf("%s skips a test", p)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// requiredKinds is the scored schema: every kind must be mutated (and caught) in at least one category.
var requiredKinds = []string{
	"steps.*.P.cols", "steps.*.P.rows", "steps.*.P.buffer", "steps.*.P.cursor.0", "steps.*.P.cursor.1", "steps.*.P.baseY",
	"steps.*.P.viewportY", "steps.*.P.normalLength", "steps.*.P.modes.*", "steps.*.P.opts.*", "steps.*.P.title",
	"steps.*.P.titles.*", "steps.*.P.bells", "steps.*.P.replies.*", "steps.*.P.lineFeeds",
	"steps.*.cells.normal.*.c.*.0", "steps.*.cells.normal.*.c.*.1", "steps.*.cells.normal.*.c.*.2", "steps.*.cells.normal.*.c.*.3",
	"steps.*.cells.normal.*.w", "steps.*.cells.alt.*.c.*.0", "steps.*.cells.alt.*.c.*.2", "steps.*.cells.alt.*.c.*.3", "steps.*.cells.alt.*.w",
	"steps.*.attrs.*.P.fg.0", "steps.*.attrs.*.P.fg.1", "steps.*.attrs.*.P.bg.0", "steps.*.attrs.*.P.bg.1", "steps.*.attrs.*.P.b",
	"steps.*.attrs.*.P.i", "steps.*.attrs.*.P.d", "steps.*.attrs.*.P.u", "steps.*.attrs.*.P.bl", "steps.*.attrs.*.P.inv",
	"steps.*.attrs.*.P.inv0", "steps.*.attrs.*.P.s", "steps.*.attrs.*.P.o",
	"B.t1.*", "B.t2.*", "B.t3", "B.t4.glyphs.*", "B.t4.pen.fg.0", "B.t4.pen.bg.0",
}
