package handler

import (
	"testing"
)

func TestNewPagination(t *testing.T) {
	pg1 := newPagination()
	if pg1.Page != 1 || pg1.Limit != 20 {
		t.FailNow()
	}
	pg2 := newPagination(2, 42)
	if pg2.Page != 2 || pg2.Limit != 42 {
		t.FailNow()
	}
}
