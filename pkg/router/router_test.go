package router

import (
	"github.com/Jarnpher553/gemini/pkg/service"
	"github.com/Jarnpher553/gemini/pkg/service/annotation"
	"testing"
)

type TestService struct {
	*service.BaseService
}

func (s *TestService) Use(handler *annotation.Handler) {
	handler.AreaN("a")
	handler.BaseRoute("wo")
}

func TestRouter_doRegister(t *testing.T) {

	r := New(Area(true))
	r.rootGroup("api")
	r.doRegister(service.NewService(&TestService{}))

}
