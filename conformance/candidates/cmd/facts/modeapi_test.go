package main

import (
	"testing"
	"testing/fstest"
)

// TestModeAPI: a mode method and an exported DecPrivateModes field are getters; test files and other files are not
// read.
func TestModeAPI(t *testing.T) {
	fsys := fstest.MapFS{
		"t.go":      {Data: []byte("package x\nfunc (t *Terminal) SetMode(p *Params) bool {\nfunc (t *Terminal) Modes() Modes { return t.m }\n")},
		"c.go":      {Data: []byte("package x\ntype CoreService struct {\n\tDecPrivateModes     DecPrivateModes\n}\n")},
		"t_test.go": {Data: []byte("package x\nfunc (t *Terminal) SetOtherMode(x int) {\n")},
		"README.md": {Data: []byte("func (t *Terminal) SetReadmeMode(x int) {\n")},
	}
	got, err := modeAPI("e", fsys)
	if want := "MODEAPI e get=DecPrivateModes,Modes replay=yes"; err != nil || got != want {
		t.Errorf("%q %v, want %q", got, err, want)
	}
	if got, err := modeAPI("e", fstest.MapFS{"a.go": {Data: []byte("package x\n")}}); err != nil || got != "MODEAPI e get=none replay=yes" {
		t.Errorf("empty: %q %v", got, err)
	}
}
