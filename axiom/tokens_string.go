// Code generated by "stringer -type=Permission -linecomment -output=tokens_string.go"; DO NOT EDIT.

package axiom

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[emptyPermission-0]
	_ = x[CanIngest-1]
	_ = x[CanQuery-2]
}

const _Permission_name = "CanIngestCanQuery"

var _Permission_index = [...]uint8{0, 0, 9, 17}

func (i Permission) String() string {
	if i >= Permission(len(_Permission_index)-1) {
		return "Permission(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Permission_name[_Permission_index[i]:_Permission_index[i+1]]
}
