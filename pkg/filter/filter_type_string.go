// Code generated by "stringer -type Type -output filter_type_string.go"; DO NOT EDIT.

package filter

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Eq-0]
	_ = x[NE-1]
	_ = x[Lt-2]
	_ = x[LtOE-3]
	_ = x[Gt-4]
	_ = x[GtOE-5]
}

const _Type_name = "EqNELtLtOEGtGtOE"

var _Type_index = [...]uint8{0, 2, 4, 6, 10, 12, 16}

func (i Type) String() string {
	if i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
