package endpoints

import (
	"context"
	"time"

	"github.com/aintsashqa/simple-todo/internal"
	"github.com/aintsashqa/simple-todo/pkg/filter"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"gopkg.in/guregu/null.v4"
)

type Endpoints struct {
	Create            endpoint.Endpoint
	Update            endpoint.Endpoint
	GetList           endpoint.Endpoint
	ChangeToCompleted endpoint.Endpoint
}

func MakeEndpoints(r internal.Repository, l log.Logger) Endpoints {
	return Endpoints{
		Create:            makeCreate(r, l),
		GetList:           makeGetList(r, l),
		Update:            makeUpdate(r, l),
		ChangeToCompleted: makeChangeToCompleted(r, l),
	}
}

func makeCreate(r internal.Repository, l log.Logger) endpoint.Endpoint {
	l = log.With(l, "endpoint", "Create")

	return func(ctx context.Context, request any) (any, error) {
		req := request.(CreateRequest)
		level.Info(l).Log("msg", "create new todo")

		// TODO: add validation
		todo := internal.Todo{
			Title:       req.Title,
			Description: req.Description,
		}

		if err := r.Create(ctx, &todo); err != nil {
			level.Error(l).Log("err", err)
			return nil, err
		}

		level.Info(l).Log("msg", "todo successfully created", "todo_id", todo.ID)
		return todo, nil
	}
}

func makeGetList(r internal.Repository, l log.Logger) endpoint.Endpoint {
	l = log.With(l, "endpoint", "GetList")

	return func(ctx context.Context, request any) (any, error) {
		req := request.(GetListRequest)
		level.Info(l).Log("msg", "get todo list", "filter_type", req.FilterType, "complited_at", req.CompletedAt)

		f := filter.Filter{
			Field: "completed_at",
			Type:  filter.Eq,
			Value: null.Time{},
		}

		if req.CompletedAt.Valid {
			f.Type = req.FilterType
			f.Value = req.CompletedAt.Time
		}

		todos, err := r.GetList(ctx, f)
		if err != nil {
			level.Error(l).Log("err", err, "filter_type", req.FilterType, "complited_at", req.CompletedAt)
			return nil, err
		}

		return todos, nil
	}
}

func makeUpdate(r internal.Repository, l log.Logger) endpoint.Endpoint {
	l = log.With(l, "endpoint", "Update")

	return func(ctx context.Context, request any) (any, error) {
		req := request.(UpdateRequest)
		level.Info(l).Log("msg", "update existing todo", "todo_id", req.ID)

		todo, err := r.GetByID(ctx, req.ID)
		if err != nil {
			level.Error(l).Log("err", err, "todo_id", req.ID)
			return nil, err
		}

		// TODO: add validation
		todo.Title = req.Title
		todo.Description = req.Description

		if err := r.Update(ctx, &todo); err != nil {
			level.Error(l).Log("err", err, "todo_id", todo.ID)
			return nil, err
		}

		level.Info(l).Log("msg", "todo successfully updated", "todo_id", todo.ID)
		return todo, nil
	}
}

func makeChangeToCompleted(r internal.Repository, l log.Logger) endpoint.Endpoint {
	l = log.With(l, "endpoint", "ChangeToCompleted")

	return func(ctx context.Context, request any) (any, error) {
		req := request.(ChangeToCompletedRequest)
		level.Info(l).Log("msg", "change complete status of todo", "todo_id", req.ID)

		todo, err := r.GetByID(ctx, req.ID)
		if err != nil {
			level.Error(l).Log("err", err, "todo_id", req.ID)
			return nil, err
		}

		todo.CompletedAt = null.TimeFrom(time.Now())
		if err := r.ChangeCompleteStatus(ctx, &todo); err != nil {
			level.Error(l).Log("err", err, "todo_id", todo.ID)
			return nil, err
		}

		level.Info(l).Log("msg", "todo complete status successfully changed", "todo_id", todo.ID)
		return todo, nil
	}
}
