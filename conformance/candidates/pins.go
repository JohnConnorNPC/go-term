//go:build pins

// Package candidates scores the surveyed Go terminal engines against the conformance goldens (issue #2). This file,
// never built, keeps every pinned candidate in go.mod while no adapter imports it yet (`go mod tidy` reads all tags).
package candidates

import (
	_ "github.com/0magnet/xterm-go/vt"
	_ "github.com/JohnConnorNPC/go-term/conformance/runner"
	_ "github.com/charmbracelet/x/vt"
	_ "github.com/gitpod-io/xterm-go"
)
