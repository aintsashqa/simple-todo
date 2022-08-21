package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/aintsashqa/simple-todo/internal/endpoints"
	"github.com/aintsashqa/simple-todo/pkg/filter"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrParseBody  = errors.New("unable parse request body to struct")
	ErrParseParam = errors.New("unable parse request uri param to variable")
	ErrParseQuery = errors.New("unable parse request query params")
)

func decodeCreateRequest(_ context.Context, r *http.Request) (any, error) {
	var result endpoints.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		return nil, ErrParseBody
	}

	return result, nil
}

func decodeUpdateRequest(_ context.Context, r *http.Request) (any, error) {
	var result endpoints.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		return nil, ErrParseBody
	}

	param := chi.URLParam(r, "todo_id")
	todoId, err := uuid.Parse(param)
	if err != nil {
		return nil, ErrParseParam
	}

	result.ID = todoId
	return result, nil
}

func decodeGetListRequest(_ context.Context, r *http.Request) (any, error) {
	const dateLayout = "02-01-2006"

	var result endpoints.GetListRequest
	query := r.URL.Query()

	if completed := query.Get("completed"); query.Has("completed") {
		values := strings.Split(completed, ":")
		if len(values) != 2 {
			return nil, ErrParseQuery
		}

		t, _ := time.Parse(dateLayout, values[1])
		result.FilterType = filter.TypeFromString(values[0])
		result.CompletedAt = null.TimeFrom(t)
	}

	return result, nil
}

func decodeChangeToCompletedRequest(_ context.Context, r *http.Request) (any, error) {
	param := chi.URLParam(r, "todo_id")
	todoId, err := uuid.Parse(param)
	if err != nil {
		return nil, ErrParseParam
	}

	return endpoints.ChangeToCompletedRequest{ID: todoId}, nil
}
