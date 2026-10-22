package summary

import "testing"

// TestJSON: sorted keys, no timestamps, the skip flag spelled out.
func TestJSON(t *testing.T) {
	s := Summary{Engine: "e", Ref: "x53", Equiv: "full", Pin: "p", GoldenManifestSha: "m", Total: 1, Passed: 0,
		Categories: map[string]Cat{"c0": {Total: 1}}, Skipped: []Skip{{ID: "c0/a", Variant: "bytes"}},
		Cases: []Case{{ID: "c0/a", Result: "SPLIT", Variant: "s0k1", Detail: "s0k1: d"}}}
	b, err := s.JSON()
	want := `{
 "cases": [
  {
   "detail": "s0k1: d",
   "id": "c0/a",
   "result": "SPLIT",
   "variant": "s0k1"
  }
 ],
 "categories": {
  "c0": {
   "passed": 0,
   "total": 1
  }
 },
 "engine": "e",
 "equiv": "full",
 "goldenManifestSha": "m",
 "passed": 0,
 "pin": "p",
 "ref": "x53",
 "skipped": [
  {
   "complete": false,
   "id": "c0/a",
   "variant": "bytes"
  }
 ],
 "total": 1
}
`
	if err != nil || string(b) != want {
		t.Fatalf("JSON %v:\n%s", err, b)
	}
}
