package httpclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Response struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	Data    *struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	} `json:"data"`
}

var s = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/test" {
		resp := Response{
			ErrCode: 200,
			ErrMsg:  "success",
			Data: &struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{Name: "lijianfeng", Age: 29},
		}

		data, _ := json.Marshal(&resp)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(200)
		writer.Write(data)
	}
}))

func TestNew(t *testing.T) {
	client := New()

	if client == nil {
		t.FailNow()
	}
}

func TestReqClient_RGet(t *testing.T) {
	t.Log(s.Listener.Addr())

	var resp Response

	client := New()
	c, _ := context.WithCancel(context.TODO())
	_ = client.RGet("http://"+s.Listener.Addr().String()+"/test", nil, c, &resp)

	t.Log(resp)
}

func TestReqClient_RPost(t *testing.T) {
	t.Log(s.Listener.Addr())

	var resp Response

	client := New()
	c, _ := context.WithCancel(context.TODO())
	_ = client.RPost("http://"+s.Listener.Addr().String()+"/test", nil, c, &resp)

	t.Log(resp)
}

func TestReqClient_RPut(t *testing.T) {
	t.Log(s.Listener.Addr())

	var resp Response

	client := New()
	c, _ := context.WithCancel(context.TODO())
	_ = client.RPut("http://"+s.Listener.Addr().String()+"/test", nil, c, &resp)

	t.Log(resp)
}

func TestReqClient_RDelete(t *testing.T) {
	t.Log(s.Listener.Addr())

	var resp Response

	client := New()
	c, _ := context.WithCancel(context.TODO())
	_ = client.RDelete("http://"+s.Listener.Addr().String()+"/test", nil, c, &resp)

	t.Log(resp)
}
