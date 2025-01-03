// Code generated by "enumer -type FileType -transform lower -output filetype_string.go"; DO NOT EDIT.

package filetype

import (
	"fmt"
	"strings"
)

const _FileTypeName = "mp3"

var _FileTypeIndex = [...]uint8{0, 3}

const _FileTypeLowerName = "mp3"

func (i FileType) String() string {
	if i >= FileType(len(_FileTypeIndex)-1) {
		return fmt.Sprintf("FileType(%d)", i)
	}
	return _FileTypeName[_FileTypeIndex[i]:_FileTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _FileTypeNoOp() {
	var x [1]struct{}
	_ = x[MP3-(0)]
}

var _FileTypeValues = []FileType{MP3}

var _FileTypeNameToValueMap = map[string]FileType{
	_FileTypeName[0:3]:      MP3,
	_FileTypeLowerName[0:3]: MP3,
}

var _FileTypeNames = []string{
	_FileTypeName[0:3],
}

// FileTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func FileTypeString(s string) (FileType, error) {
	if val, ok := _FileTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _FileTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to FileType values", s)
}

// FileTypeValues returns all values of the enum
func FileTypeValues() []FileType {
	return _FileTypeValues
}

// FileTypeStrings returns a slice of all String values of the enum
func FileTypeStrings() []string {
	strs := make([]string, len(_FileTypeNames))
	copy(strs, _FileTypeNames)
	return strs
}

// IsAFileType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i FileType) IsAFileType() bool {
	for _, v := range _FileTypeValues {
		if i == v {
			return true
		}
	}
	return false
}
