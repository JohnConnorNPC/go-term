package summary

import "encoding/json"

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

// JSON is the summary's bytes: no timestamps, one trailing newline.
func (s Summary) JSON() ([]byte, error) {
	b, err := json.MarshalIndent(s, "", " ")
	return append(b, '\n'), err
}
