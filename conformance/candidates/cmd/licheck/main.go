// Command licheck prints every LIC line of the three candidates (atoms 25-27), read from the module cache's verified
// extraction of each pinned zip (internal/modcache):
//
//	go run ./cmd/licheck > scores/lic.txt
//
// Exit 2 when a module or its imports can't be read.
package main

import (
	"fmt"
	"os"
)

// engines are the candidates and their modules; gitpod's lines gate (E8).
var engines = [][2]string{{"gitpod", "github.com/gitpod-io/xterm-go"}, {"magnet", "github.com/0magnet/xterm-go"},
	{"xvt", "github.com/charmbracelet/x/vt"}}

func main() {
	code := 0
	for _, e := range engines {
		lines, err := check(e[0] == "gitpod", e[1])
		if err != nil {
			fmt.Println("LIC", e[0], "L0 FAIL unreadable:", err)
			code = 2
		}
		for _, l := range lines {
			fmt.Println("LIC", e[0], l.ID, l.Status, l.Value)
		}
	}
	os.Exit(code)
}
