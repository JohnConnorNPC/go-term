// Package scoring is cmd/score's body: the candidate registry, the run over every golden, and the summary file.
package scoring

import (
	"runtime/debug"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/gitpod"
	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/safe"
	"github.com/JohnConnorNPC/go-term/conformance/candidates/magnet"
	"github.com/JohnConnorNPC/go-term/conformance/candidates/xvt"
	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// candidates are the scored engines by name, with the module each one adapts (its version is the summary's pin).
var candidates = map[string]struct {
	new    runner.Factory
	module string
}{
	"gitpod": {gitpod.New, "github.com/gitpod-io/xterm-go"},
	"magnet": {magnet.New, "github.com/0magnet/xterm-go"},
	"xvt":    {xvt.New, "github.com/charmbracelet/x/vt"},
}

// Lookup returns the named engine: a candidate behind safe (E4), or one of the runner's anchors, stub and oracle,
// unwrapped (safe would hide the oracle's private golden hook).
func Lookup(name string, pr *runner.Probes) (runner.Factory, bool) {
	if c, ok := candidates[name]; ok {
		return safe.Wrap(c.new, safe.Default), true
	}
	if name == "oracle" {
		return runner.Oracle(pr), true
	}
	return runner.Lookup(name)
}

// Pin is the candidate's module@version from the build info, "-" for an anchor.
func Pin(name string) string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "-"
	}
	for _, d := range info.Deps {
		if d.Path == candidates[name].module {
			return d.Path + "@" + d.Version
		}
	}
	return "-"
}
