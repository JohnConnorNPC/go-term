package screen

import (
	"encoding/hex"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/blake2b"
)

// Cursor is a cursor position, 0-based within the grid (never counting
// history rows).
type Cursor struct{ Row, Col int }

// HashSize is the digest size in bytes: blake2b with digest_size=16, which is
// a different hash from blake2b-512 truncated to 16 bytes.
const HashSize = 16

// ContentHash is Studio's content_hash: blake2b-128 of the UTF-8 text, as 32
// lowercase hex digits. Invalid UTF-8 is replaced by U+FFFD first (Python's
// errors="replace" writes "?" for a lone surrogate, which a Go string from vt
// never holds).
func ContentHash(text string) string {
	h, err := blake2b.New(HashSize, nil)
	if err != nil {
		panic(err) // only for a size outside 1..64 or a key over 64 bytes
	}
	h.Write([]byte(strings.ToValidUTF8(text, "�")))
	return hex.EncodeToString(h.Sum(nil))
}

// StableHash is Studio's stable_hash, as Python's _stable_hash: the
// ContentHash of text with the cursor's character replaced by NUL. The
// character is at row c.Row+rowOffset of text split on "\n" (rowOffset is the
// HistoryLines prepended on a scrollback read) and code point (rune) index
// c.Col of that row. A nil cursor or one outside the text gives ContentHash.
//
// Read does not call this with a cell column: it maps the cursor cell to the
// runes of its glyph itself (a wide glyph is 2 cells and 1 rune, a combined
// cell is 1 cell and several runes); see Frame.StableHash.
func StableHash(text string, c *Cursor, rowOffset int) string {
	if c == nil {
		return ContentHash(text)
	}
	rows := strings.Split(text, "\n")
	r := c.Row + rowOffset
	if r < 0 || r >= len(rows) || c.Col < 0 {
		return ContentHash(text)
	}
	start := runeOffset(rows[r], c.Col)
	if start < 0 {
		return ContentHash(text)
	}
	_, size := utf8.DecodeRuneInString(rows[r][start:])
	return ContentHash(maskRow(rows, r, start, start+size))
}

// runeOffset is the byte offset of rune i of s, or -1 when s has no rune i.
func runeOffset(s string, i int) int {
	n := 0
	for off := range s {
		if n == i {
			return off
		}
		n++
	}
	return -1
}

// maskRow joins rows with "\n" after replacing bytes [start,end) of row r by
// one NUL.
func maskRow(rows []string, r, start, end int) string {
	var b strings.Builder
	for i, row := range rows {
		if i > 0 {
			b.WriteByte('\n')
		}
		if i == r {
			b.WriteString(row[:start])
			b.WriteByte(0)
			b.WriteString(row[end:])
			continue
		}
		b.WriteString(row)
	}
	return b.String()
}
