package canonicalized

import (
	"fmt"
	"strconv"
)

func (v *Version) compareType(other *Version) bool {
	return fetchType(v) == fetchType(other)
}

func fetchType(v *Version) string {
	typeTag := ""
	for _, _type := range v.Type {
		typeTag = fmt.Sprintf(
			"%s.%s.%s",
			typeTag,
			_type.Name,
			strconv.FormatInt(
				_type.Tag,
				10),
		)
	}

	return typeTag
}

func isSudo(v *Version) bool {
	if len(v.Timestamp) > 0 && len(v.CommitHash) > 0 {
		return true
	}

	return false
}

func isBeta(v *Version) bool {
	if isStable(v) {
		return false
	}

	for _, t := range v.Type {
		if t.Name == "beta" {
			return true
		}
	}

	return false
}

func isAlpha(v *Version) bool {
	if isStable(v) {
		return false
	}

	for _, t := range v.Type {
		if t.Name == "alpha" {
			return true
		}
	}

	return false
}

func isStable(v *Version) bool {
	if len(v.Type) == 0 {
		return true
	}

	for _, t := range v.Type {
		if t.Name == "stable" {
			return true
		}
	}

	return false
}

func toInt64PtrFromString(s string) *int64 {
	if s == "" {
		return nil
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil
	}
	return &n
}
