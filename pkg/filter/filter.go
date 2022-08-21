//go:generate stringer -type Type -output filter_type_string.go
package filter

import "strings"

const (
	Eq Type = iota
	NE
	Lt
	LtOE
	Gt
	GtOE
)

type Type uint8

type Filter struct {
	Field string
	Type  Type
	Value any
}

func TypeFromString(v string) Type {
	switch strings.ToLower(v) {
	case strings.ToLower(NE.String()):
		return NE
	case strings.ToLower(Lt.String()):
		return Lt
	case strings.ToLower(LtOE.String()):
		return LtOE
	case strings.ToLower(Gt.String()):
		return Gt
	case strings.ToLower(GtOE.String()):
		return GtOE
	}

	return Eq
}
