// Command facts prints the report-only facts of the three candidates: GOLINE and MODEAPI per engine, then the
// GOLDEN-VOCAB lines (atoms 28, 28b):
//
//	go run ./cmd/facts > scores/facts.txt
//
// Exit 2 when a module or the typings can't be read. The typings are #1's build output, ../gen/node_modules (not in a
// fresh checkout: conformance/scripts/fetch-pins.sh on the host, then npm ci --offline in node.sh); the tests read a
// vendored excerpt (cmd/facts/testdata) instead, which scripts/vocab-check.sh ties to them.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/modcache"
)

var engines = [][2]string{{"gitpod", "github.com/gitpod-io/xterm-go"}, {"magnet", "github.com/0magnet/xterm-go"},
	{"xvt", "github.com/charmbracelet/x/vt"}}

func main() {
	code := 0
	for _, e := range engines {
		dir, err := modcache.Dir(e[1])
		gomod, err2 := os.ReadFile(filepath.Join(dir, "go.mod"))
		api, err3 := modeAPI(e[0], os.DirFS(dir))
		if err != nil || err2 != nil || err3 != nil {
			fmt.Println("FACTS", e[0], "UNREADABLE", err, err2, err3)
			code = 2
			continue
		}
		fmt.Println(goLine(e[0], gomod))
		fmt.Println(api)
	}
	vocab, err := goldenVocab("../gen/node_modules/xterm-headless/typings/xterm-headless.d.ts")
	if err != nil {
		fmt.Println("FACTS vocab UNREADABLE", err)
		code = 2
	}
	for _, l := range vocab {
		fmt.Println(l)
	}
	os.Exit(code)
}
