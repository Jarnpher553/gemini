package jwt

import "testing"

func TestNew(t *testing.T) {
	_, err := New(1)
	if err != nil {
		t.FailNow()
	}
}

func TestParse(t *testing.T) {
	data, _ := New(1)

	_, err := Parse(data)
	if err != nil {
		t.FailNow()
	}
}
