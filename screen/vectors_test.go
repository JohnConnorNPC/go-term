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
// browserland-studio vectors/ into testdata/: the first six at commit b6053e0,
// raw, the headless colours, pyte and raw replays and headless/divergence.json
// (Studio's list of where headless and Python differ) at ed3249b, where the
// first six are unchanged. The sniffer families (altscreen, synthetic, titles)
// are read from ../sniff/testdata, pinned here as there. The headless files'
// headers carry the sha256 of the python file each was replayed from.
var vectorPins = map[string]string{
	"python/raw.json":                     "f11e8e445d3f31f9f21cd77c89d2751607538c32a5e197fa4c11de1af6458c5b",
	"headless/raw.json":                   "f20f3d19299c4ed9df9a6550ec3df6f12184d90adc5de3738dd9b4edfc568899",
	"headless/colours.json":               "76bb882e1b93f4a5c004228f32df539fd8a9132142ae1e2405867600be76bc94",
	"headless/pyte.json":                  "59e8b52ddd5b2548714a785b4eb4b58b7c2f383b1fa6d5534ba7dc288364f10a",
	"headless/divergence.json":            "34c726244eed10d4f15852abd64e5da472f94ef265e2a091a524e16117e594dc",
	sniffData + "python/altscreen.json":   "c5221a92bec87e0a2791380b9d836c019888977aa17833b5510d264cc9c397ed",
	sniffData + "python/synthetic.json":   "d3ef7b87b05939878688cb0ce80dc6ef49eb757f44b7c71c7ca254a028f86803",
	sniffData + "python/titles.json":      "58c0384d54a2e57622a0e96f557d8908e83b8abb22ad29e0e90ce8c3c2cce791",
	sniffData + "headless/altscreen.json": "6a60401833fddd30145a54dce28e7173b5c70447459d4b2805d77ee6c61d0612",
	sniffData + "headless/synthetic.json": "d5407c330659f18390c2d328c8592287ded4bc39acc8ca3253fbda9ee5634046",
	sniffData + "headless/titles.json":    "63cb5a8595e91d3d20d941a34a9e6d908fdcf0e068ea2e1c26a1bfbd8e061879",
	"python/colours.json":                 "6c4638852b1fecbd58decc47742cfd545db50b61232c0764d6439a086593266f",
	"python/hashes.json":                  "cc87245017305fdcf6005f6da0df0e77e590ecb5996427542fab587a7573e1df",
	"python/pyte.json":                    "7a4e88416b06d85090c9be5084cc1dd2ab15bbb74540f04bc040682364814fce",
	"python/textgrid.json":                "349ab3f09955c1978feeec4aad9a2d76ed316ffe0878972c5777bad372512dbd",
	"headless/hashes.json":                "b71eb62c0b898757ef24b320fe99f043ccc839e692c7aa2decca281e334a2f69",
	"headless/textgrid.json":              "ff63f85b45da7a0dcf80b6101cba5018a8fc76b6e467e58635ebb8c9e2097ad5",
}

// sniffData is the sniff package's vector dir, relative to testdata/.
const sniffData = "../../sniff/testdata/"

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
	// Ops are a stream case's [method, args, result, state after] (the sniffers).
	Ops [][4]json.RawMessage `json:"ops"`
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
