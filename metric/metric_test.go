package metric

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	printer := NewPrinter()

	m := New(printer, time.Second*1)

	if m == nil {
		t.FailNow()
	}
}

func TestMetric_Start(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.FailNow()
		}
	}()

	printer := NewPrinter()

	m := New(printer, time.Second*1)

	m.Stop()

	m.Start()
}

func TestMetric_Stop(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.FailNow()
		}
	}()

	printer := NewPrinter()

	m := New(printer, time.Second*1)

	m.Stop()
}
