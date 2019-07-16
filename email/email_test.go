package email

import "testing"

func TestBind(t *testing.T) {
	Bind(Host("bjmail.hylinkad.com"), Port(25), Username("magellan@hylinkad.com"), Pwd("hy8888"))

	err := Send("测试", "哈哈", "349030388@qq.com")

	if err != nil {
		t.FailNow()
	}
}
