package service

type Handler struct {
	Middleware   []Middleware
	RelativePath string
}

func (h *Handler) UseMiddleware(m ...Middleware) {
	h.Middleware = append(h.Middleware, m...)
}

func (h *Handler) Path(path string) {
	h.RelativePath = path
}
