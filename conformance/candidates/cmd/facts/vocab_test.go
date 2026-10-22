package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// excerpt is the vendored excerpt of the 5.3.0 typings (its header names the whole file's sha256), so the test needs
// no gen/node_modules.
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

// TestVocabSource: the excerpt's lines are lines of the real typings, whose sha256 is the one the excerpt names. Without
// gen/node_modules (a fresh checkout: fetch-pins.sh + npm ci make it) this is an explicit SKIP, counted by the suite.
func TestVocabSource(t *testing.T) {
	typings, err := os.ReadFile("../../../gen/node_modules/xterm-headless/typings/xterm-headless.d.ts")
	if os.IsNotExist(err) {
		t.Skip("SKIP vocab-source: ../gen/node_modules absent")
	}
	ex, err2 := os.ReadFile(excerpt)
	if err != nil || err2 != nil {
		t.Fatal(err, err2)
	}
	sum := sha256.Sum256(typings)
	m := regexp.MustCompile(`sha256 ([0-9a-f]{64})`).FindSubmatch(ex)
	if m == nil || string(m[1]) != hex.EncodeToString(sum[:]) {
		t.Errorf("excerpt names %q, the typings are %x", m, sum)
	}
	lines := strings.Split(string(typings), "\n")
	for _, l := range strings.Split(strings.TrimSpace(string(ex)), "\n") {
		if !strings.HasPrefix(l, "//") && !slices.Contains(lines, l) {
			t.Errorf("excerpt line %q is not in the typings", l)
		}
	}
}
