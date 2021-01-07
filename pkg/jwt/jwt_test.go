package jwt

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tokenStr, err := New(map[string]interface{}{
		"demo": 1,
	}, time.Minute*5)
	if err != nil {
		t.FailNow()
	}
	t.Log(tokenStr)
}

func TestParse(t *testing.T) {
	data, _ := New(map[string]interface{}{
		"demo": 1,
	}, time.Minute*5)

	token, err := Parse(data)
	if err != nil {
		t.FailNow()
	}
	t.Log(token)
}
