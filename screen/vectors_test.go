package screen

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// vectorPins are the sha256 of Studio's vector files copied byte for byte from
// browserland-studio vectors/ (commit b6053e0) into testdata/. The headless
// files' headers carry the sha256 of the python file each was replayed from.
var vectorPins = map[string]string{
	"python/hashes.json":     "cc87245017305fdcf6005f6da0df0e77e590ecb5996427542fab587a7573e1df",
	"python/textgrid.json":   "349ab3f09955c1978feeec4aad9a2d76ed316ffe0878972c5777bad372512dbd",
	"headless/hashes.json":   "b71eb62c0b898757ef24b320fe99f043ccc839e692c7aa2decca281e334a2f69",
	"headless/textgrid.json": "ff63f85b45da7a0dcf80b6101cba5018a8fc76b6e467e58635ebb8c9e2097ad5",
}

func readVectors(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestVectorPins(t *testing.T) {
	names := make([]string, 0, len(vectorPins))
	for n := range vectorPins {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		sum := sha256.Sum256(readVectors(t, n))
		if got := hex.EncodeToString(sum[:]); got != vectorPins[n] {
			t.Errorf("%s: sha256 %s, pinned %s", n, got, vectorPins[n])
		}
	}
}

// pyCase is one case of a python vector file. Args and dict outputs are
// {"$d": [[key, value], ...]}, Python's ordered dict encoding.
type pyCase struct {
	ID     string          `json:"id"`
	Entry  string          `json:"entry"`
	Args   json.RawMessage `json:"args"`
	Output json.RawMessage `json:"output"`
	Replay json.RawMessage `json:"replay"` // {chunks_b64, cols, rows, ...}; headless input_sha256 is its sha256
}

func loadPython(t *testing.T, name string) []pyCase {
	t.Helper()
	var f struct {
		Cases []pyCase `json:"cases"`
	}
	if err := json.Unmarshal(readVectors(t, name), &f); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	if len(f.Cases) == 0 {
		t.Fatalf("%s: no cases (unmeasured)", name)
	}
	return f.Cases
}

// pyDict decodes {"$d": [[key, value], ...]} into a map of raw values; ok is
// false when raw is not such an object (null, for example).
func pyDict(t *testing.T, raw json.RawMessage) (m map[string]json.RawMessage, ok bool) {
	t.Helper()
	var d struct {
		D *[][2]json.RawMessage `json:"$d"`
	}
	if json.Unmarshal(raw, &d) != nil || d.D == nil {
		return nil, false
	}
	m = map[string]json.RawMessage{}
	for _, kv := range *d.D {
		var k string
		if err := json.Unmarshal(kv[0], &k); err != nil {
			t.Fatalf("dict key %s: %v", kv[0], err)
		}
		m[k] = kv[1]
	}
	return m, true
}

func mustUnmarshal(t *testing.T, raw json.RawMessage, v any) {
	t.Helper()
	if err := json.Unmarshal(raw, v); err != nil {
		t.Fatalf("%s: %v", raw, err)
	}
}
