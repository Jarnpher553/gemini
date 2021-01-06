package email

import (
	"crypto/tls"
	"fmt"
	"github.com/go-gomail/gomail"
	"regexp"
)

var dialer *gomail.Dialer

type Option func(map[string]interface{})

func Host(host string) Option {
	return func(m map[string]interface{}) {
		m["host"] = host
	}
}

func Port(port int) Option {
	return func(m map[string]interface{}) {
		m["port"] = port
	}
}

func Username(name string) Option {
	return func(m map[string]interface{}) {
		m["username"] = name
	}
}

func Pwd(pwd string) Option {
	return func(m map[string]interface{}) {
		m["pwd"] = pwd
	}
}

func Bind(opts ...Option) {
	m := make(map[string]interface{})

	for _, opt := range opts {
		opt(m)
	}

	dialer = gomail.NewDialer(m["host"].(string), m["port"].(int), m["username"].(string), m["pwd"].(string))
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
}

func Send(subject string, content string, to ...string) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	msg := gomail.NewMessage()
	msg.SetHeader("From", dialer.Username)
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", subject)

	reg := regexp.MustCompile(`<.*>.*</.*>`)
	match := reg.MatchString(content)
	if match {
		msg.SetBody("text/html", content)
	} else {
		msg.SetBody("text/plain", content)
	}

	return dialer.DialAndSend(msg)
}
