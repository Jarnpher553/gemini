package erro

import "testing"

func TestRegister(t *testing.T) {
	Register(10000, "hahaha")

	t.Log(ErrMsg[10000])

	if ErrMsg[10000] != "hahaha" {
		t.FailNow()
	}
}
