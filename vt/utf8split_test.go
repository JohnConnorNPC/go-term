package vt

import "testing"

// TestUtf8SplitEveryPartition [#14]: a character whose held bytes include a 0x80 continuation byte decodes the same
// when the writes cut it anywhere, in every partition of its bytes into writes (2^(n-1) cut sets); malformed input
// does too. xterm.js's Utf8ToUtf32, and @xterm/headless 6.0.0 with it, lose the character on a write ending after a
// held 0x80 (its goldens wide/cjk-ext-b, wide/zwj-family, wide/zwsp are split.sensitive): vt differs there on purpose.
func TestUtf8SplitEveryPartition(t *testing.T) {
	streams := []string{"aĀb", "aࠀb", "aကb", "a—b", "a​b", "a‍b", "aꀀb",
		"a\U00010000b", "a\U00020000b", "a\U00040000b", "a\U00100000b", "aé€\U0001f389b",
		"a\xe2\x80Xb", "a\xf0\xa0\x80Xb", "a\xf0\xa0Xb", "a\xe2\x80\x80\x80b", "a\xc0\x80b", "a\xe0\x80\x80b",
		"a\xed\xa0\x80b", "a\xf4\x90\x80\x80b", "a\xef\xbb\xbfb", "a\xff\x80b", "a\xe2\x80\xe2\x80\x94b"}
	decode := func(b []byte, mask int) string { // a cut before byte k when bit k-1 of mask is set
		var dec Utf8ToUtf32
		var out []rune
		target := make([]uint32, len(b)+4)
		for start, k := 0, 1; k <= len(b); k++ {
			if k == len(b) || mask&(1<<(k-1)) != 0 {
				for _, cp := range target[:dec.Decode(b[start:k], target)] {
					out = append(out, rune(cp)) // #nosec G115 -- decoder output is at most 0x10FFFF
				}
				start = k
			}
		}
		return string(out)
	}
	parts := 0
	for _, s := range streams {
		whole := decode([]byte(s), 0)
		for mask := 1; mask < 1<<(len(s)-1); mask++ {
			parts++
			if got := decode([]byte(s), mask); got != whole {
				t.Errorf("%q cut mask %b: got %q, whole %q", s, mask, got, whole)
			}
		}
	}
	t.Logf("UTF8SPLIT streams=%d partitions=%d", len(streams), parts)
}
