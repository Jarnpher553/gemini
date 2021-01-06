package router

import (
	"github.com/Jarnpher553/gemini/pkg/service"
	"testing"
)

type TestService struct {
	*service.BaseService
}

func (s *TestService) Use(handler *service.Handler) {
	handler.AreaN("a")
	handler.BaseRoute("wo")
}

func TestRouter_doRegister(t *testing.T) {

	r := New(Area(true))
	r.rootGroup("api")
	r.doRegister(service.NewService(&TestService{}))

}
