package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aintsashqa/simple-todo/internal"
)

type response struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func encodeResponse(_ context.Context, w http.ResponseWriter, v any) error {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{Data: v})
	return nil
}

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	statusCode := http.StatusInternalServerError

	switch err {
	case internal.ErrNotFound:
		statusCode = http.StatusNotFound
		break
	case ErrParseBody, ErrParseParam, ErrParseQuery:
		statusCode = http.StatusBadRequest
		break
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response{Message: err.Error()})
}
