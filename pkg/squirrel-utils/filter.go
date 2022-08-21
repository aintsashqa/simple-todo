package squirrelutils

import (
	"github.com/Masterminds/squirrel"
	"github.com/aintsashqa/simple-todo/pkg/filter"
)

type Filter filter.Filter

func (f Filter) Condition() squirrel.Sqlizer {
	switch f.Type {
	case filter.NE:
		return squirrel.NotEq{f.Field: f.Value}
	case filter.Lt:
		return squirrel.Lt{f.Field: f.Value}
	case filter.LtOE:
		return squirrel.LtOrEq{f.Field: f.Value}
	case filter.Gt:
		return squirrel.Gt{f.Field: f.Value}
	case filter.GtOE:
		return squirrel.GtOrEq{f.Field: f.Value}
	case filter.Eq:
	default:
	}

	return squirrel.Eq{f.Field: f.Value}
}
