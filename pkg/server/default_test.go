package server

import "testing"

func TestDefault(t *testing.T) {
	s := Default()
	if s == nil {
		t.FailNow()
	}
}
