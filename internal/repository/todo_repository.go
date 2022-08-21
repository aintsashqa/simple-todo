package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/aintsashqa/simple-todo/internal"
	"github.com/aintsashqa/simple-todo/pkg/entity"
	"github.com/aintsashqa/simple-todo/pkg/filter"
	squirrelutils "github.com/aintsashqa/simple-todo/pkg/squirrel-utils"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

const (
	table string = schema + "todo"

	c_ID          string = "id"
	c_Title       string = "title"
	c_Description string = "description"
	c_CreatedAt   string = "created_at"
	c_UpdatedAt   string = "updated_at"
	c_DeletedAt   string = "deleted_at"
	c_CompletedAt string = "completed_at"
)

var (
	ErrBuildSqlQuery      = errors.New("unable to build sql query")
	ErrExecSqlQuery       = errors.New("unable to execute sql query")
	ErrScanResult         = errors.New("unable scan result to struct")
	ErrExecOrScanSqlQuery = fmt.Errorf("%s or %s", ErrExecSqlQuery, ErrScanResult)
)

type repo struct {
	c Client
	b squirrel.StatementBuilderType
	l log.Logger
}

func NewRepository(c Client, l log.Logger) *repo {
	return &repo{
		c: c,
		b: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		l: log.With(l, "repository", "todo", "driver", "postgresql"),
	}
}

func (r *repo) exec(ctx context.Context, sql string, args []any) error {
	l := log.With(r.l, "method", "exec")

	level.Debug(l).Log("sql", sql, "args", args)
	if tag, err := r.c.Exec(ctx, sql, args...); err != nil {
		if tag.RowsAffected() == 0 {
			level.Error(l).Log("err", internal.ErrNotFound, "cause", err)
			return internal.ErrNotFound
		}

		level.Error(l).Log("err", ErrExecSqlQuery, "cause", err)
		return ErrExecSqlQuery
	}

	return nil
}

func (r *repo) Create(ctx context.Context, todo *internal.Todo) error {
	l := log.With(r.l, "method", "Create")
	todo.Before()

	sql, args, err := r.b.
		Insert(table).
		Columns(c_ID, c_Title, c_Description, c_CreatedAt, c_UpdatedAt).
		Values(todo.ID, todo.Title, todo.Description, todo.CreatedAt, todo.UpdatedAt).
		ToSql()
	if err != nil {
		level.Error(l).Log("err", ErrBuildSqlQuery, "cause", err)
		return ErrBuildSqlQuery
	}

	return r.exec(ctx, sql, args)
}

func (r *repo) Update(ctx context.Context, todo *internal.Todo) error {
	l := log.With(r.l, "method", "Update")
	todo.Before()

	sql, args, err := r.b.
		Update(table).
		SetMap(map[string]any{
			c_Title:       todo.Title,
			c_Description: todo.Description,
			c_UpdatedAt:   todo.UpdatedAt,
		}).
		Where(squirrel.Eq{
			c_ID:        todo.ID,
			c_DeletedAt: nil,
		}).
		ToSql()
	if err != nil {
		level.Error(l).Log("err", ErrBuildSqlQuery, "cause", err)
		return ErrBuildSqlQuery
	}

	return r.exec(ctx, sql, args)
}

func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (internal.Todo, error) {
	l := log.With(r.l, "method", "GetByID")

	sql, args, err := r.b.
		Select(c_Title, c_Description, c_CreatedAt, c_UpdatedAt).
		From(table).
		Where(squirrel.Eq{
			c_ID:          id,
			c_DeletedAt:   nil,
			c_CompletedAt: nil,
		}).
		ToSql()
	if err != nil {
		level.Error(l).Log("err", ErrBuildSqlQuery, "cause", err)
		return internal.Todo{}, ErrBuildSqlQuery
	}

	level.Debug(l).Log("sql", sql, "args", args)
	todo := internal.Todo{Base: entity.Base{ID: id}}
	if err := r.c.QueryRow(ctx, sql, args...).
		Scan(&todo.Title, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			level.Error(l).Log("err", internal.ErrNotFound, "cause", err)
			return internal.Todo{}, internal.ErrNotFound
		}

		level.Error(l).Log("err", ErrExecOrScanSqlQuery, "cause", err)
		return internal.Todo{}, ErrExecOrScanSqlQuery
	}

	return todo, nil
}

func (r *repo) GetList(ctx context.Context, filters ...filter.Filter) ([]internal.Todo, error) {
	l := log.With(r.l, "method", "GetList")

	b := r.b.
		Select(c_ID, c_Title, c_Description, c_CreatedAt, c_UpdatedAt, c_DeletedAt, c_CompletedAt).
		From(table)

	for _, filter := range filters {
		b = b.Where(squirrelutils.Filter(filter).Condition())
	}

	sql, args, err := b.ToSql()
	if err != nil {
		level.Error(l).Log("err", ErrBuildSqlQuery, "cause", err)
		return nil, ErrBuildSqlQuery
	}

	level.Debug(l).Log("sql", sql, "args", args)
	records, err := r.c.Query(ctx, sql, args...)
	if err != nil {
		level.Error(l).Log("err", ErrExecSqlQuery, "cause", err)
		return nil, ErrExecSqlQuery
	}

	var todos []internal.Todo
	for records.Next() {
		var todo internal.Todo
		if err := records.Scan(
			&todo.ID, &todo.Title, &todo.Description,
			&todo.CreatedAt, &todo.UpdatedAt, &todo.DeletedAt, &todo.CompletedAt,
		); err != nil {
			level.Error(l).Log("err", ErrScanResult, "cause", err)
			return nil, ErrScanResult
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *repo) ChangeCompleteStatus(ctx context.Context, todo *internal.Todo) error {
	l := log.With(r.l, "method", "ChangeCompleteStatus")
	todo.Before()

	sql, args, err := r.b.
		Update(table).
		SetMap(map[string]any{
			c_UpdatedAt:   todo.UpdatedAt,
			c_CompletedAt: todo.CompletedAt,
		}).
		Where(squirrel.Eq{
			c_ID:          todo.ID,
			c_DeletedAt:   nil,
			c_CompletedAt: nil,
		}).
		ToSql()
	if err != nil {
		level.Error(l).Log("err", ErrBuildSqlQuery, "cause", err)
		return ErrBuildSqlQuery
	}

	return r.exec(ctx, sql, args)
}
