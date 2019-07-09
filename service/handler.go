package service

type Handler struct {
	Middleware []Middleware
}

func (h *Handler) UseMiddleware(m ...Middleware) {
	h.Middleware = append(h.Middleware, m...)
}
