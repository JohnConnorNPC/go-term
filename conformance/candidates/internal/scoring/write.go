package scoring

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/count"
	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/summary"
	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// Write checks the denominator (internal/count), builds the summary of rs over the required categories and writes
// it to path. Exit 0 every case passes, 1 a case fails, 2 UNMEASURED (no summary, no file).
func Write(out io.Writer, id summary.Summary, root, path string, rs []runner.Result) int {
	if err := count.Check(root, rs); err != nil {
		fmt.Fprintln(out, err)
		return 2
	}
	var cats []string
	for c := range runner.RequiredCategories {
		cats = append(cats, c)
	}
	sort.Strings(cats)
	man, err := os.ReadFile(filepath.Join(root, "golden", "manifest.json"))
	sum := sha256.Sum256(man)
	id.GoldenManifestSha = hex.EncodeToString(sum[:])
	s, berr := summary.Build(id, rs, cats)
	b, jerr := s.JSON()
	if err != nil || berr != nil || jerr != nil {
		fmt.Fprintln(out, "UNMEASURED summary", err, berr, jerr)
		return 2
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		fmt.Fprintln(out, "UNMEASURED write", err)
		return 2
	}
	fmt.Fprintf(out, "TOTAL %d/%d\n", s.Passed, s.Total)
	if s.Passed != s.Total {
		return 1
	}
	return 0
}
