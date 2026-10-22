package summary

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// Summary is the E6 score file. Fields are declared in key order, so the JSON keys come out sorted.
type Summary struct {
	Cases             []Case         `json:"cases"`
	Categories        map[string]Cat `json:"categories"`
	Engine            string         `json:"engine"`
	Equiv             string         `json:"equiv"`
	GoldenManifestSha string         `json:"goldenManifestSha"`
	Passed            int            `json:"passed"`
	Pin               string         `json:"pin"`
	Ref               string         `json:"ref"`
	Skipped           []Skip         `json:"skipped"`
	Total             int            `json:"total"`
}

// Case is one scored case; Detail is the runner's Diff verbatim.
type Case struct {
	Detail  string `json:"detail"`
	ID      string `json:"id"`
	Result  string `json:"result"`
	Variant string `json:"variant"`
}

// Cat is one category's totals.
type Cat struct {
	Passed int `json:"passed"`
	Total  int `json:"total"`
}

// Skip is one variant headless itself is sensitive to. Complete is true for a PASS case; a SPLIT-FAIL case's list
// stops at its first failing variant, so it is only a lower bound (rulings-5 #5).
type Skip struct {
	Complete bool   `json:"complete"`
	ID       string `json:"id"`
	Variant  string `json:"variant"`
}

// Build fills s (whose identity fields are set) from one run's results; every required category must have results.
func Build(s Summary, results []runner.Result, categories []string) (Summary, error) {
	s.Categories, s.Cases, s.Skipped = map[string]Cat{}, []Case{}, []Skip{}
	for _, c := range categories {
		s.Categories[c] = Cat{}
	}
	for _, r := range results {
		cat, known := s.Categories[r.Category]
		if !known {
			return s, fmt.Errorf("%s: category %s is not required", r.ID, r.Category)
		}
		class, v := Classify(r)
		cat.Total++
		if class == "PASS" {
			cat.Passed++
		}
		s.Categories[r.Category] = cat
		s.Cases = append(s.Cases, Case{Detail: r.Diff, ID: r.ID, Result: class, Variant: v})
		for _, sk := range r.Skipped {
			s.Skipped = append(s.Skipped, Skip{Complete: class == "PASS", ID: r.ID, Variant: sk})
		}
	}
	s.Total, s.Passed = 0, 0
	for _, c := range categories {
		if s.Categories[c].Total == 0 {
			return s, fmt.Errorf("category %s has no results", c)
		}
		s.Total, s.Passed = s.Total+s.Categories[c].Total, s.Passed+s.Categories[c].Passed
	}
	sort.Slice(s.Cases, func(i, j int) bool { return s.Cases[i].ID < s.Cases[j].ID })
	sort.SliceStable(s.Skipped, func(i, j int) bool { return s.Skipped[i].ID < s.Skipped[j].ID })
	return s, nil
}

// JSON is the summary's bytes: no timestamps, one trailing newline.
func (s Summary) JSON() ([]byte, error) {
	b, err := json.MarshalIndent(s, "", " ")
	return append(b, '\n'), err
}
