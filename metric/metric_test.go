package metric

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	writer := NewWriter(time.Second * 2)

	m := New(writer)

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

	writer := NewWriter(time.Second * 2)

	m := New(writer)

	m.Stop()

	m.Start()
}

func TestMetric_Stop(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.FailNow()
		}
	}()

	writer := NewWriter(time.Second * 2)

	m := New(writer)

	m.Stop()
}
