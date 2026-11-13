package runner

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Status is the outcome of one scored case.
type Status string

const (
	PASS  Status = "PASS"
	FAIL  Status = "FAIL"  // a whole-write step or a probe differs
	STALE Status = "STALE" // the golden's ops no longer hash to its stored inSha
	SPLIT Status = "SPLIT" // whole writes pass, a split variant differs
)

// Result is one case's score; Diff is the first mismatch. Skipped names the split variants not scored because headless
// itself is sensitive to them (S2): they are reported, never a silent PASS. Rescored names those same variants when
// Options.ScoreSensitive scored them instead.
type Result struct {
	ID, Category string
	Status       Status
	Diff         string
	Skipped      []string
	Rescored     []string
}

// Options select what RunCase scores beyond the whole-write run and the probes. Splits scores the chunk-split
// variants; ScoreSensitive (the -fuzz mode) also scores those headless itself is sensitive to, against the
// whole-write golden as every other variant: go-term is split-safe at every cut (#14), so any difference is SPLIT.
type Options struct{ Splits, ScoreSensitive bool }

const splitLimit = 512 // as gen/lib/split.mjs SPLIT_LIMIT

func start(f Factory, g *Golden) (Engine, error) {
	e, err := f(Config{g.Init.Cols, g.Init.Rows, g.Init.Scrollback, g.Init.Unicode})
	if ca, ok := e.(caseAware); ok && err == nil {
		ca.begin(g)
	}
	return e, err
}

func apply(e Engine, op Op, chunks [][]byte) error {
	switch op.Op {
	case "w":
		if chunks == nil {
			chunks = [][]byte{op.Bytes}
		}
		for _, c := range chunks {
			if err := e.Write(c); err != nil {
				return err
			}
		}
		return nil
	case "resize":
		return e.Resize(op.Cols, op.Rows)
	default:
		return e.Reset()
	}
}

type variant struct {
	id     string
	chunks func(i int, b []byte) [][]byte
}

// variants mirrors gen/lib/split.mjs: "bytes" (one byte per write), then "s<i>k<k>" two-chunk splits of each write
// step i when the case writes fewer than splitLimit bytes in total.
func variants(g *Golden) []variant {
	v := []variant{{"bytes", func(_ int, b []byte) [][]byte {
		out := [][]byte{{}}[:0]
		for i := range b {
			out = append(out, b[i:i+1])
		}
		if len(out) == 0 {
			out = [][]byte{{}}
		}
		return out
	}}}
	total := 0
	for _, op := range g.Ops {
		total += len(op.Bytes)
	}
	for i, op := range g.Ops {
		for k := 1; total < splitLimit && op.Op == "w" && k < len(op.Bytes); k++ {
			i, k := i, k
			v = append(v, variant{fmt.Sprintf("s%dk%d", i, k), func(j int, b []byte) [][]byte {
				if j != i {
					return nil
				}
				return [][]byte{b[:k], b[k:]}
			}})
		}
	}
	return v
}

// steps feeds every op (chunked by chunk, nil = whole) and compares the snapshot after each with the golden.
func steps(e Engine, g *Golden, chunk func(int, []byte) [][]byte) string {
	for i, op := range g.Ops {
		var c [][]byte
		if chunk != nil && op.Op == "w" {
			c = chunk(i, op.Bytes)
		}
		if err := apply(e, op, c); err != nil {
			return fmt.Sprintf("step[%d]: engine error: %v", i, err)
		}
		st, err := e.Snapshot()
		if err != nil {
			return fmt.Sprintf("step[%d]: snapshot error: %v", i, err)
		}
		if d := Compare(fmt.Sprintf("step[%d]", i), g.Steps[i], st); d != "" {
			return d
		}
	}
	return ""
}

// RunCase scores one golden: STALE check, the whole-write steps, the probe twins, then (opt.Splits) every split
// variant except those headless itself flagged as chunk-sensitive, unless opt.ScoreSensitive.
func RunCase(g *Golden, f Factory, pr *Probes, opt Options) Result {
	r := Result{ID: g.ID, Category: g.Category(), Status: PASS}
	fail := func(s Status, d string) Result { r.Status, r.Diff = s, d; return r }
	if g.Stale() {
		return fail(STALE, g.stale+" inSha "+g.InSha+" ops "+InSha(g.Ops))
	}
	e, err := start(f, g)
	if err != nil {
		return fail(FAIL, "engine: "+err.Error())
	}
	if d := steps(e, g, nil); d != "" {
		return fail(FAIL, d)
	}
	b, err := probe(f, g, pr)
	var u unobservedErr
	if errors.As(err, &u) {
		return fail(FAIL, u.Error())
	}
	if err != nil {
		return fail(FAIL, "B: "+err.Error())
	}
	if d := CompareB(g.B, b); d != "" {
		return fail(FAIL, d)
	}
	skip := map[string]bool{}
	for _, id := range g.Sensitive {
		skip[id] = true
	}
	for _, v := range variants(g) {
		if !opt.Splits {
			continue
		}
		if skip[v.id] && !opt.ScoreSensitive {
			r.Skipped = append(r.Skipped, v.id)
			continue
		}
		if skip[v.id] {
			r.Rescored = append(r.Rescored, v.id)
		}
		if e, err = start(f, g); err != nil {
			return fail(FAIL, "engine: "+err.Error())
		}
		if d := steps(e, g, v.chunks); d != "" {
			return fail(SPLIT, v.id+": "+d)
		}
	}
	return r
}

func fileSha(p string) string {
	b, err := os.ReadFile(p)
	if err != nil {
		return ""
	}
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}

// LoadAll parses every golden of one ref under root (the conformance directory), sorted by id, and marks the ones whose
// provenance does not match the tree (S3).
func LoadAll(root, ref string) ([]*Golden, error) {
	var out []*Golden
	var man struct{ Cases map[string]string }
	if b, err := os.ReadFile(filepath.Join(root, "golden", "manifest.json")); err != nil || json.Unmarshal(b, &man) != nil {
		return nil, fmt.Errorf("manifest: %v", err)
	}
	pins := fileSha(filepath.Join(root, "gen", "pins.json"))
	err := filepath.WalkDir(filepath.Join(root, "golden", ref), func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".json") {
			return err
		}
		g, err := Load(p)
		if err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
		g.stale = provenance(g, ref, fileSha(filepath.Join(root, "cases", g.ID+".json")), man.Cases[g.ID], pins)
		out = append(out, g)
		return nil
	})
	return out, err
}
