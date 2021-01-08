package annotation

import (
	"strings"
)

type Handler struct {
	Middleware   []interface{}
	RelativePath string
	HttpMethod   string
	BasePath     string
	UseArea      bool
	AreaName     string
}

type Middleware func(middleware []interface{}) []interface{}

func (h *Handler) Use(ms ...Middleware) {
	for _, m := range ms {
		h.Middleware = m(h.Middleware)
	}
}

func (h *Handler) Route(httpMethod string, path string) {
	h.HttpMethod = strings.ToTitle(httpMethod)
	h.RelativePath = path
}

func (h *Handler) Post(path string) {
	h.Route("POST", path)
}

func (h *Handler) Get(path string) {
	h.Route("GET", path)
}

func (h *Handler) Put(path string) {
	h.Route("PUT", path)
}

func (h *Handler) Delete(path string) {
	h.Route("DELETE", path)
}

func (h *Handler) Options(path string) {
	h.Route("OPTIONS", path)
}

func (h *Handler) Patch(path string) {
	h.Route("PATCH", path)
}

func (h *Handler) Head(path string) {
	h.Route("HEAD", path)
}

func (h *Handler) BaseRoute(path string) {
	h.BasePath = path
}

func (h *Handler) Area(use bool) {
	h.UseArea = use
}

func (h *Handler) AreaN(name string) {
	h.AreaName = name
}
