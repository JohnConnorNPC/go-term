package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestDeterministic: the same seed writes byte-identical files; another seed writes different ones.
func TestDeterministic(t *testing.T) {
	src := filepath.Join("..", "..")
	a, b, c := t.TempDir(), t.TempDir(), t.TempDir()
	for _, run := range []struct {
		dir  string
		seed uint64
	}{{a, 7}, {b, 7}, {c, 8}} {
		if _, err := write(run.dir, src, run.seed, 20); err != nil {
			t.Fatal(err)
		}
	}
	same, differ := 0, 0
	for i := range 20 {
		fa, _ := os.ReadFile(filepath.Join(a, "cases", "fuzz", name(7, i)))
		fb, _ := os.ReadFile(filepath.Join(b, "cases", "fuzz", name(7, i)))
		fc, _ := os.ReadFile(filepath.Join(c, "cases", "fuzz", name(8, i)))
		if len(fa) == 0 || !bytes.Equal(fa, fb) {
			t.Fatalf("seed 7 case %d: runs differ or empty (%d bytes)", i, len(fa))
		}
		same++
		if !bytes.Equal(fa[bytes.IndexByte(fa, ','):], fc[bytes.IndexByte(fc, ','):]) {
			differ++
		}
	}
	if same != 20 || differ < 19 {
		t.Fatalf("same=%d differ=%d, want 20 and >= 19", same, differ)
	}
	for _, f := range []string{"probes.json", "gen/pins.json"} {
		want, _ := os.ReadFile(filepath.Join(src, f))
		got, err := os.ReadFile(filepath.Join(a, f))
		if err != nil || len(want) == 0 || !bytes.Equal(got, want) {
			t.Fatalf("%s not copied byte for byte (err %v)", f, err)
		}
	}
}

func name(seed uint64, i int) string {
	c := Generate(seed, i)
	return filepath.Base(c.ID) + ".json"
}

// TestBounds: 2000 cases decode as the node loader decodes them (exactly one key per step, even-length hex, resize in
// bounds, reset true) and hold the bounds: geometry, steps, scrollback, noop-ok, most under the 512-byte split limit.
func TestBounds(t *testing.T) {
	under, kinds := 0, map[string]int{}
	const n = 2000
	for i := range n {
		c := Generate(42, i)
		var back map[string]any
		b, _ := json.Marshal(c)
		if err := json.Unmarshal(b, &back); err != nil {
			t.Fatal(err)
		}
		keys := []string{}
		for k := range back {
			keys = append(keys, k)
		}
		if len(keys) != 6 || c.Scrollback != Scrollback || !reflect.DeepEqual(c.Tags, []string{"noop-ok"}) {
			t.Fatalf("%s: keys %v scrollback %d tags %v", c.ID, keys, c.Scrollback, c.Tags)
		}
		if c.Cols < MinCols || c.Cols > MaxCols || c.Rows < MinRows || c.Rows > MaxRows {
			t.Fatalf("%s: size %dx%d", c.ID, c.Cols, c.Rows)
		}
		if len(c.Steps) < 1 || len(c.Steps) > MaxSteps {
			t.Fatalf("%s: %d steps", c.ID, len(c.Steps))
		}
		for _, st := range c.Steps {
			if len(st) != 1 {
				t.Fatalf("%s: step %v", c.ID, st)
			}
			switch {
			case st["x"] != nil:
				x, err := hex.DecodeString(st["x"].(string))
				if err != nil || len(x) == 0 || len(x) > MaxString {
					t.Fatalf("%s: x %q: %v", c.ID, st["x"], err)
				}
				kinds["x"]++
			case st["resize"] != nil:
				r := st["resize"].([]int)
				if r[0] < MinCols || r[0] > MaxCols || r[1] < MinRows || r[1] > MaxRows {
					t.Fatalf("%s: resize %v", c.ID, r)
				}
				kinds["resize"]++
			case st["reset"] == true:
				kinds["reset"]++
			default:
				t.Fatalf("%s: unknown step %v", c.ID, st)
			}
		}
		if c.writeBytes() < 512 {
			under++
		}
	}
	if under < n*95/100 || kinds["x"] == 0 || kinds["resize"] == 0 || kinds["reset"] == 0 {
		t.Fatalf("under 512 bytes: %d/%d (want >= 95%%), step kinds %v", under, n, kinds)
	}
	t.Logf("FUZZGEN cases=%d under512=%d steps %v", n, under, kinds)
}
