package runner

import (
	"encoding/json"
	"fmt"
	"sort"
)

func show(v any) string { b, _ := json.Marshal(v); return string(b) }

func mismatch(path string, want, got any) string {
	return fmt.Sprintf("%s: want %s got %s", path, show(want), show(got))
}

// sameStrings treats nil and empty as equal.
func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func compareMap(path string, want, got map[string]any) string {
	keys := map[string]bool{}
	for k := range want {
		keys[k] = true
	}
	for k := range got {
		keys[k] = true
	}
	var ks []string
	for k := range keys {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	for _, k := range ks {
		w, wok := want[k]
		g, gok := got[k]
		if wok != gok || show(w) != show(g) {
			return mismatch(path+"."+k, w, g)
		}
	}
	return ""
}

// Compare returns "" when got matches want on every scored (tier P) field, else the first mismatching path.
func Compare(prefix string, want, got State) string {
	w, g := want.P, got.P
	ints := []struct {
		name string
		w, g int
	}{{"cols", w.Cols, g.Cols}, {"rows", w.Rows, g.Rows}, {"cursor.x", w.Cursor[0], g.Cursor[0]}, {"cursor.y", w.Cursor[1], g.Cursor[1]},
		{"baseY", w.BaseY, g.BaseY}, {"viewportY", w.ViewportY, g.ViewportY}, {"normalLength", w.NormalLength, g.NormalLength},
		{"bells", w.Bells, g.Bells}, {"lineFeeds", w.LineFeeds, g.LineFeeds}}
	for _, f := range ints {
		if f.w != f.g {
			return mismatch(prefix+".P."+f.name, f.w, f.g)
		}
	}
	switch {
	case w.Buffer != g.Buffer:
		return mismatch(prefix+".P.buffer", w.Buffer, g.Buffer)
	case w.Title != g.Title:
		return mismatch(prefix+".P.title", w.Title, g.Title)
	case !sameStrings(w.Titles, g.Titles):
		return mismatch(prefix+".P.titles", w.Titles, g.Titles)
	case !sameStrings(w.Replies, g.Replies):
		return mismatch(prefix+".P.replies", w.Replies, g.Replies)
	}
	if d := compareMap(prefix+".P.modes", w.Modes, g.Modes); d != "" {
		return d
	}
	if d := compareMap(prefix+".P.opts", w.Opts, g.Opts); d != "" {
		return d
	}
	if d := compareRows(prefix+".normal", want.Normal, got.Normal); d != "" {
		return d
	}
	if (want.Alt == nil) != (got.Alt == nil) {
		return mismatch(prefix+".alt", want.Alt != nil, got.Alt != nil)
	}
	return compareRows(prefix+".alt", want.Alt, got.Alt)
}

func compareRows(path string, want, got []Row) string {
	if len(want) != len(got) {
		return mismatch(path+".len", len(want), len(got))
	}
	for y := range want {
		if want[y].Wrapped != got[y].Wrapped {
			return mismatch(fmt.Sprintf("%s[%d].wrapped", path, y), want[y].Wrapped, got[y].Wrapped)
		}
		wc, gc := want[y].Cells, got[y].Cells
		if len(wc) != len(gc) {
			return mismatch(fmt.Sprintf("%s[%d].len", path, y), len(wc), len(gc))
		}
		for x := range wc {
			p := fmt.Sprintf("%s[%d][%d]", path, y, x)
			switch {
			case wc[x].Ch != gc[x].Ch:
				return mismatch(p+".ch", wc[x].Ch, gc[x].Ch)
			case wc[x].W != gc[x].W:
				return mismatch(p+".w", wc[x].W, gc[x].W)
			case wc[x].A != gc[x].A:
				return mismatch(p+".attr", wc[x].A, gc[x].A)
			}
		}
	}
	return ""
}

// CompareB compares the raw probe results (tier B).
func CompareB(want, got B) string {
	if d := compareMap("B.t1", toAny(want.T1), toAny(got.T1)); d != "" {
		return d
	}
	switch {
	case show(want.T2) != show(got.T2):
		return mismatch("B.t2", want.T2, got.T2)
	case want.T3 != got.T3:
		return mismatch("B.t3", want.T3, got.T3)
	case want.T4.Pen != got.T4.Pen:
		return mismatch("B.t4.pen", want.T4.Pen, got.T4.Pen)
	case !sameStrings(want.T4.Glyphs, got.T4.Glyphs):
		return mismatch("B.t4.glyphs", want.T4.Glyphs, got.T4.Glyphs)
	}
	return ""
}

func toAny(m map[string]string) map[string]any {
	out := map[string]any{}
	for k, v := range m {
		out[k] = v
	}
	return out
}
