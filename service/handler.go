package service

import "strings"

type Handler struct {
	Middleware   []Middleware
	RelativePath string
	HttpMethod   string
}

func (h *Handler) UseMiddleware(m ...Middleware) {
	h.Middleware = append(h.Middleware, m...)
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
