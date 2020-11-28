package router

import (
	"github.com/Jarnpher553/gemini/service"
	"testing"
)

type TestService struct {
	*service.BaseService
}

func (s *TestService) Area() string {
	return ""
}

func TestRouter_doRegister(t *testing.T) {

	r := New(Area(true))
	r.rootGroup("api")
	r.doRegister(service.NewService(&TestService{}))

}
