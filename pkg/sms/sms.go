package sms

import (
	"fmt"
	"github.com/Jarnpher553/gemini/pkg/httpclient"
)

var sms Sms

type Sms interface {
	AuthToken(client *httpclient.ReqClient) (string, error)
	Send(client *httpclient.ReqClient, token string, smsParam map[string]string, phone ...string) error
}

func Bind(s Sms) {
	sms = s
}

func Send(client *httpclient.ReqClient, token string, smsParam map[string]string, phone ...string) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()
	return sms.Send(client, token, smsParam, phone...)
}

func AuthToken(client *httpclient.ReqClient) (s string, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()
	return sms.AuthToken(client)
}
