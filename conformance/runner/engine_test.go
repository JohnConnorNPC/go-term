package runner

import "testing"

func TestEngineStub(t *testing.T) {
	f, ok := Lookup("stub")
	if !ok {
		t.Fatalf("stub not registered: %v", Engines())
	}
	e, _ := f(Config{Cols: 20, Rows: 6})
	if st, err := e.Snapshot(); err != nil || st.P.Cols != 0 || st.Normal != nil {
		t.Fatalf("stub snapshot not empty: %+v %v", st, err)
	}
}
