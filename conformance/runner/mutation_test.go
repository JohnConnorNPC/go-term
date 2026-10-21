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

// pinned is the oracle bound to the pristine golden: Begin is shadowed, so it never sees the mutated expected side.
type pinned struct{ *oracle }

func (pinned) Begin(*Golden) {}

var (
	index    = regexp.MustCompile(`\.\d+(\.|$)`)
	freeKeys = regexp.MustCompile(`^(steps\.\*\.P\.(modes|opts)|B\.t1)\.[^.]+$`)
	// Not scored on purpose: identity/pins (STALE has its own test), display text, tier I, split bookkeeping.
	unscored = regexp.MustCompile(`^(id|ref|pins|inSha|init|split)(\.|$)|^steps\.\*\.(op|text|I)(\.|$)|^steps\.\*\.attrs\.\*\.I(\.|$)`)
)

func kindOf(path string) string {
	k := path
	for index.MatchString(k) {
		k = index.ReplaceAllString(k, ".*$1")
	}
	if m := freeKeys.FindStringSubmatch(k); m != nil {
		k = m[1] + ".*"
	}
	return k
}

// leaves calls f for every scalar leaf (not null) with its dotted path; f may replace the leaf's value.
func leaves(v any, path string, f func(path string, v any) (any, bool)) any {
	switch t := v.(type) {
	case map[string]any:
		for k, x := range t {
			t[k] = leaves(x, strings.TrimPrefix(path+"."+k, "."), f)
		}
	case []any:
		for i, x := range t {
			t[i] = leaves(x, path+"."+itoa(i), f)
		}
	case nil:
	default:
		if nv, ok := f(path, v); ok {
			return nv
		}
	}
	return v
}

func itoa(i int) string { b, _ := json.Marshal(i); return string(b) }

func mutate(path string, v any) any {
	switch t := v.(type) {
	case bool:
		return !t
	case string:
		return t + "Z"
	case json.Number:
		if strings.HasSuffix(path, ".1") && strings.Contains(path, ".cells.") {
			return json.Number("999") // an attr index: point outside the table (two entries may share their P part)
		}
		n, _ := t.Int64()
		return json.Number(itoa(int(n) + 1))
	}
	return v
}

func decode(t *testing.T, b []byte) map[string]any {
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	var m map[string]any
	if err := d.Decode(&m); err != nil {
		t.Fatal(err)
	}
	return m
}

func TestMutationExactlyOne(t *testing.T) {
	pr, err := LoadProbes("../probes.json")
	if err != nil {
		t.Fatal(err)
	}
	files, _ := filepath.Glob("../golden/x53/*/*.json")
	cats := map[string]bool{}
	first := map[string][2]string{} // kind -> (file, path)
	raw := map[string][]byte{}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		raw[f] = b
		cats[filepath.Base(filepath.Dir(f))] = true
		g, err := Parse(b)
		if err != nil {
			t.Fatal(err)
		}
		if r := RunCase(g, Oracle(pr), pr, Options{}); r.Status != PASS { // control arm: every pristine case passes
			t.Fatalf("pristine %s: %s %s", g.ID, r.Status, r.Diff)
		}
		leaves(decode(t, b), "", func(p string, v any) (any, bool) {
			if k := kindOf(p); !unscored.MatchString(k) {
				if _, seen := first[k]; !seen {
					first[k] = [2]string{f, p}
				}
			}
			return nil, false
		})
	}
	var want []string
	for c := range RequiredCategories {
		want = append(want, c)
	}
	var got []string
	for c := range cats {
		got = append(got, c)
	}
	sort.Strings(want)
	sort.Strings(got)
	if len(got) == 0 || strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("golden categories %v, want %v", got, want)
	}
	for _, k := range requiredKinds {
		if _, ok := first[k]; !ok {
			t.Errorf("leaf kind %s appears in no golden", k)
		}
	}
	var kinds []string
	for k := range first {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)
	for _, k := range kinds {
		f, p := first[k][0], first[k][1]
		pristine, _ := Parse(raw[f])
		m := decode(t, raw[f])
		leaves(m, "", func(q string, v any) (any, bool) { return mutate(q, v), q == p })
		mb, _ := json.Marshal(m)
		want, err := Parse(mb)
		if err != nil { // a mutation the strict parser rejects is a failed case too
			continue
		}
		fac := func(Config) (Engine, error) { return pinned{&oracle{pr: pr, step: -1, g: pristine}}, nil }
		if r := RunCase(want, fac, pr, Options{}); r.Status == PASS {
			t.Errorf("kind %s (%s in %s): 0 cases failed, want 1", k, p, pristine.ID)
		}
	}
	t.Logf("kinds=%d categories=%d cases=%d", len(kinds), len(cats), len(files))
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

// requiredKinds is the scored schema: every one must occur in some golden, and mutating it must fail its case.
var requiredKinds = []string{
	"steps.*.P.cols", "steps.*.P.rows", "steps.*.P.buffer", "steps.*.P.cursor.*", "steps.*.P.baseY", "steps.*.P.viewportY",
	"steps.*.P.normalLength", "steps.*.P.modes.*", "steps.*.P.opts.*", "steps.*.P.title", "steps.*.P.titles.*", "steps.*.P.bells",
	"steps.*.P.replies.*", "steps.*.P.lineFeeds",
	"steps.*.cells.normal.*.c.*.*", "steps.*.cells.normal.*.w", "steps.*.cells.alt.*.c.*.*", "steps.*.cells.alt.*.w",
	"steps.*.attrs.*.P.fg.*", "steps.*.attrs.*.P.bg.*", "steps.*.attrs.*.P.b", "steps.*.attrs.*.P.i", "steps.*.attrs.*.P.d",
	"steps.*.attrs.*.P.u", "steps.*.attrs.*.P.bl", "steps.*.attrs.*.P.inv", "steps.*.attrs.*.P.inv0", "steps.*.attrs.*.P.s", "steps.*.attrs.*.P.o",
	"B.t1.*", "B.t2.*", "B.t3", "B.t4.glyphs.*", "B.t4.pen.fg.*", "B.t4.pen.bg.*",
}
