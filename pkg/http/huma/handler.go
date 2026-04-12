package huma

import (
	"net/http"
)

// Handler abstraction for the router to be used with Huma.
type Handler struct {
	config         *Config
	router         Router
	operations     []OperationWithHandlerRegisterer
	middleware     []MiddlewareFunc
	formatter      []Formatter
	securityScheme []SecurityScheme
}

// ServeHTTP an http handler.
func (h *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	h.router.ServeHTTP(writer, req)
}
