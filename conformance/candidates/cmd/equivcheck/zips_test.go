package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fixture writes a file proxy holding module example.com/m v1.0.0 with files, and a go.mod that requires it.
func fixture(t *testing.T, files map[string]string) (proxy, gomod string) {
	t.Helper()
	dir := t.TempDir()
	proxy, gomod = filepath.Join(dir, "proxy"), filepath.Join(dir, "go.mod")
	if err := os.MkdirAll(filepath.Join(proxy, "example.com/m/@v"), 0o755); err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(filepath.Join(proxy, "example.com/m/@v/v1.0.0.zip"))
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(f)
	for name, body := range files {
		fw, err := w.Create("example.com/m@v1.0.0/" + name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := fw.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(gomod, []byte("module x\n\nrequire example.com/m v1.0.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return proxy, gomod
}

// TestModule: the zip's files are read by their path inside the module, at the go.mod version; a module without a
// require line, or with an upper-case path, is an error.
func TestModule(t *testing.T) {
	proxy, gomod := fixture(t, map[string]string{"vt/a.go": "package vt\nx := 1\n"})
	z := &zips{proxy: proxy, gomod: gomod, mods: map[string]map[string][]string{}}
	files, err := z.module("example.com/m")
	if err != nil || strings.Join(files["vt/a.go"], "|") != "package vt|x := 1|" {
		t.Fatalf("vt/a.go: %q %v", files["vt/a.go"], err)
	}
	for _, mod := range []string{"example.com/other", "example.com/M"} {
		if _, err := z.module(mod); err == nil {
			t.Errorf("module %s: no error", mod)
		}
	}
}
