package main

import "testing"

// TestGoldenVocab: the real 5.3.0 typings give the two unions the table maps to.
func TestGoldenVocab(t *testing.T) {
	p := "../../../gen/node_modules/xterm-headless/typings/xterm-headless.d.ts"
	got, err := goldenVocab(p)
	if err != nil || len(got) != 2 || got[0] != "GOLDEN-VOCAB mouseTrackingMode none|x10|vt200|drag|any ("+p+")" ||
		got[1] != "GOLDEN-VOCAB cursorStyle block|underline|bar ("+p+")" {
		t.Errorf("%q %v", got, err)
	}
	if _, err := goldenVocab("nope.d.ts"); err == nil {
		t.Error("missing typings accepted")
	}
}
