package screen

import (
	"strings"
	"testing"
)

// TestBlake2b128Pin: the digest is blake2b with digest_size=16, Python's
// hashlib.blake2b(b"", digest_size=16), not blake2b-512 cut to 16 bytes
// (786a02f742015903c6c6fd852552d272).
func TestBlake2b128Pin(t *testing.T) {
	if got := ContentHash(""); got != "cae66941d9efbd404e4d88758ea67670" {
		t.Errorf("ContentHash(\"\") = %s", got)
	}
	if got := ContentHash(""); strings.HasPrefix(got, "786a02f7") {
		t.Errorf("ContentHash is blake2b-512 truncated: %s", got)
	}
	if got := ContentHash("ab\ncd"); len(got) != 32 || strings.ToLower(got) != got {
		t.Errorf("ContentHash = %q, want 32 lowercase hex digits", got)
	}
}

func TestStableHash(t *testing.T) {
	const text = "$ ls  \nfile  "
	c := &Cursor{Row: 0, Col: 2}
	if ContentHash(text) != ContentHash("$ ls  \nfile  ") {
		t.Error("equal text, different hashes")
	}
	// A blinking cursor toggles the cell under it: content changes, stable does not.
	blink := "$ #s  \nfile  "
	if ContentHash(blink) == ContentHash(text) {
		t.Error("blink: content hash did not change")
	}
	if StableHash(blink, c, 0) != StableHash(text, c, 0) {
		t.Error("blink in place changed the stable hash")
	}
	// A move exposes another cell: the stable hash changes.
	if StableHash(text, &Cursor{Row: 0, Col: 3}, 0) == StableHash(text, c, 0) {
		t.Error("cursor move kept the stable hash")
	}
	if StableHash(text, c, 0) == ContentHash(text) {
		t.Error("an in-bounds cursor did not mask")
	}
	for _, oob := range []*Cursor{nil, {Row: 2, Col: 0}, {Row: 0, Col: 6}, {Row: -1, Col: 0}, {Row: 0, Col: -1}} {
		if StableHash(text, oob, 0) != ContentHash(text) {
			t.Errorf("cursor %+v: stable != content", oob)
		}
	}
	// rowOffset: the cursor is grid-relative, history rows sit above it.
	if StableHash("hist\n"+text, c, 1) != ContentHash("hist\n$ \x00s  \nfile  ") {
		t.Error("rowOffset 1 did not shift the mask")
	}
	// Col is a rune index: the mask replaces one whole code point.
	if StableHash("中é!", &Cursor{Row: 0, Col: 1}, 0) != ContentHash("中\x00!") {
		t.Error("rune index masking")
	}
}

// TestPythonHashVectors: Studio's python/hashes.json _content_hash and
// _stable_hash cases give byte-identical values through ContentHash and
// StableHash (the definition is exact; SA1 Q2). Python's malformed cursors
// ({} or a missing key) fall back to the content hash, as nil does here.
func TestPythonHashVectors(t *testing.T) {
	var content, contentOK, stable, stableOK int
	for _, c := range loadPython(t, "python/hashes.json") {
		var want string
		switch c.Entry {
		case "hashes._content_hash", "hashes._stable_hash":
			mustUnmarshal(t, c.Output, &want)
		default:
			continue
		}
		args, ok := pyDict(t, c.Args)
		if !ok {
			t.Fatalf("%s: args not a dict", c.ID)
		}
		var text string
		mustUnmarshal(t, args["text"], &text)
		var got string
		if c.Entry == "hashes._content_hash" {
			content++
			got = ContentHash(text)
		} else {
			stable++
			var off int
			mustUnmarshal(t, args["row_offset"], &off)
			var cur *Cursor
			if d, ok := pyDict(t, args["cursor"]); ok && d["row"] != nil && d["col"] != nil {
				cur = &Cursor{}
				mustUnmarshal(t, d["row"], &cur.Row)
				mustUnmarshal(t, d["col"], &cur.Col)
			}
			got = StableHash(text, cur, off)
		}
		if got != want {
			t.Errorf("%s %s: %s, Python %s", c.ID, c.Entry, got, want)
			continue
		}
		if c.Entry == "hashes._content_hash" {
			contentOK++
		} else {
			stableOK++
		}
	}
	if content != 69 || stable != 68 {
		t.Errorf("read %d content and %d stable cases, want 69 and 68", content, stable)
	}
	t.Logf("VECTORS python hashes content %d/%d stable %d/%d", contentOK, content, stableOK, stable)
}
