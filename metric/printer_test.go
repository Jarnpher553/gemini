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

	printer := NewPrinter()

	if printer == nil {
		t.FailNow()
	}

}

func TestLogWriter_Write(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.FailNow()
		}
	}()

	printer := NewPrinter()

	if printer == nil {
		t.FailNow()
	}

	m := New(&Config{"", printer, time.Second * 1})

	m.Stop()

	m.Start()
}
