package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/kstsm/wb-comment-tree/internal/service"
	"net/http"
)

type HandlerI interface {
	NewRouter() http.Handler
}

type Handler struct {
	service service.ServiceI
}

func NewHandler(service service.ServiceI) HandlerI {
	return &Handler{
		service: service,
	}
}

func (h *Handler) NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("web/")))

	r.Route("/comments", func(r chi.Router) {
		r.Post("/", h.createCommentHandler)
		r.Get("/", h.getCommentsHandler)
		r.Delete("/{id}", h.deleteCommentHandler)
	})

	return r
}
