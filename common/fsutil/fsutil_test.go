package fsutil

import "testing"

func TestEnsureDir(t *testing.T) {
	hiDir := "testdata/hi"
	if err := EnsureDir(hiDir); err != nil {
		t.FailNow()
	}
	if err := EnsureDir(hiDir + "/test"); err != nil {
		t.FailNow()
	}
}
