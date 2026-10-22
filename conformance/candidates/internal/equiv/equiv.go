// Package equiv is the closed equivalence table (E10): the only place an adapter may turn a candidate value that the
// golden spells differently into the golden's value. Every row is data in equiv/table.json, each one cited in the
// candidate's source (cmd/equivcheck proves the citations). The table is chosen once per process by Load.
package equiv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Row maps one enumerated candidate key to a golden value: a literal, or an @accessor the adapter reads through.
// Cite holds the row's citations; cmd/equivcheck parses and proves them.
type Row struct {
	ID, Engine, Field, Key, Class string
	Golden                        any
	Cite                          []json.RawMessage
}

// Accessors is the closed set of read-through golden values.
var Accessors = map[string]bool{"@override": true, "@Options.CursorStyle": true, "@Options.CursorBlink": true}

var active []Row

// StrictLabel goes on every strict score line (go-v4 item 4, rulings-6): strict has no table, so it measures how far
// the raw encodings are from the golden's, the mapping cost, not the engine's quality.
const StrictLabel = "strict=raw-encoding (mapping cost, not a quality score)"

// Label is " "+StrictLabel for the strict mode, else "": the suffix of a score line.
func Label(mode string) string {
	if mode == "strict" {
		return " " + StrictLabel
	}
	return ""
}

// Read parses a table file strictly: unknown keys, an unknown class or an accessor outside Accessors are errors.
func Read(path string) ([]Row, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var f struct {
		Rows []Row `json:"rows"`
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&f); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	for _, r := range f.Rows {
		s, isString := r.Golden.(string)
		switch {
		case r.Class != "case" && r.Class != "default" && r.Class != "source":
			return nil, fmt.Errorf("%s: row %s class %q", path, r.ID, r.Class)
		case isString && strings.HasPrefix(s, "@") && !Accessors[s]:
			return nil, fmt.Errorf("%s: row %s accessor %s", path, r.ID, s)
		}
	}
	return f.Rows, nil
}

// Load makes the rows of path that mode keeps the active table: strict none, case/default only that class, full all.
func Load(path, mode string) error {
	rows, err := Read(path)
	if err != nil {
		return err
	}
	if mode != "strict" && mode != "case" && mode != "default" && mode != "full" {
		return fmt.Errorf("equiv mode %q is not strict|case|default|full", mode)
	}
	active = nil
	for _, r := range rows {
		if mode == "full" || r.Class == mode {
			active = append(active, r)
		}
	}
	return nil
}

// Map returns the golden value of the active row for (engine, field, key); ok is false when no row applies, and
// the adapter then reports the candidate's raw value.
func Map(engine, field, key string) (golden any, ok bool) {
	for _, r := range active {
		if r.Engine == engine && r.Field == field && r.Key == key {
			return r.Golden, true
		}
	}
	return nil, false
}
