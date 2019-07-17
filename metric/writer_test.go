package metric

import (
	"testing"
	"time"
)

func TestNewWriter(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.FailNow()
		}
	}()

	writer := NewWriter(time.Second * 2)

	if writer == nil {
		t.FailNow()
	}

}

func TestLogWriter_Write(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.FailNow()
		}
	}()

	writer := NewWriter(time.Second * 2)

	if writer == nil {
		t.FailNow()
	}

	m := New(writer)

	m.Stop()

	m.Start()
}
