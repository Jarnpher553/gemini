package scheduler

import "testing"

func TestBind(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.FailNow()
		}
	}()

	Bind()
}
