// Code generated by "stringer -type=JobSelectionDataLocality -output=input_locality_type_string.go"; DO NOT EDIT.

package semantic

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Local-0]
	_ = x[Anywhere-1]
}

const _JobSelectionDataLocality_name = "LocalAnywhere"

var _JobSelectionDataLocality_index = [...]uint8{0, 5, 13}

func (i JobSelectionDataLocality) String() string {
	if i < 0 || i >= JobSelectionDataLocality(len(_JobSelectionDataLocality_index)-1) {
		return "JobSelectionDataLocality(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _JobSelectionDataLocality_name[_JobSelectionDataLocality_index[i]:_JobSelectionDataLocality_index[i+1]]
}