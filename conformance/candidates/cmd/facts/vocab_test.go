package main

import "testing"

// excerpt is the vendored excerpt of the 5.3.0 typings (its header names the whole file's sha256), so the test needs
// no gen/node_modules; scripts/vocab-check.sh ties it to the real typings.
const excerpt = "testdata/xterm-headless-5.3.0.d.ts"

// TestGoldenVocab: the 5.3.0 typings give the two unions the table maps to.
func TestGoldenVocab(t *testing.T) {
	got, err := goldenVocab(excerpt)
	if err != nil || len(got) != 2 || got[0] != "GOLDEN-VOCAB mouseTrackingMode none|x10|vt200|drag|any ("+excerpt+")" ||
		got[1] != "GOLDEN-VOCAB cursorStyle block|underline|bar ("+excerpt+")" {
		t.Errorf("%q %v", got, err)
	}
	if _, err := goldenVocab("nope.d.ts"); err == nil {
		t.Error("missing typings accepted")
	}
}
