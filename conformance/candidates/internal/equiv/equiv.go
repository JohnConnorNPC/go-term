// Package equiv is the closed equivalence table (E10): the only place an adapter may turn a candidate value that the
// golden spells differently into the golden's value. Every row is data in equiv/table.json, each one cited in the
// candidate's source (cmd/equivcheck proves the citations).
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
