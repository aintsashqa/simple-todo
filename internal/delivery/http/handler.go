package http

import (
	"net/http"

	"github.com/aintsashqa/simple-todo/internal"
	"github.com/aintsashqa/simple-todo/internal/endpoints"
	"github.com/go-chi/chi/v5"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHandler(r internal.Repository, l log.Logger) http.Handler {
	h := chi.NewRouter()

	endpoint := endpoints.MakeEndpoints(r, l)
	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeErrorResponse),
	}

	h.Route("/api/todo", func(h chi.Router) {
		h.Post("/", httptransport.NewServer(
			endpoint.Create,
			decodeCreateRequest,
			encodeResponse,
			opts...,
		).ServeHTTP)

		h.Get("/", httptransport.NewServer(
			endpoint.GetList,
			decodeGetListRequest,
			encodeResponse,
			opts...,
		).ServeHTTP)

		h.Put("/{todo_id}", httptransport.NewServer(
			endpoint.Update,
			decodeUpdateRequest,
			encodeResponse,
			opts...,
		).ServeHTTP)

		h.Post("/{todo_id}", httptransport.NewServer(
			endpoint.ChangeToCompleted,
			decodeChangeToCompletedRequest,
			encodeResponse,
			opts...,
		).ServeHTTP)
	})

	return h
}
