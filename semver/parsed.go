package semver

import (
	"github.com/rng70/versions/canonicalized"
)

func NewVersionFromList(versions []string) []*canonicalized.Version {
	var version []*canonicalized.Version

	for _, s := range versions {
		p := NewVersion(s)

		version = append(version, &p)
	}

	return version
}

func NewVersion(s string) canonicalized.Version {
	p := canonicalized.ParseVersionString(s)

	return p
}

func SortedParsedVersions(versions []string, descending ...bool) []*canonicalized.Version {
	out := NewVersionFromList(versions)

	desc := false
	if len(descending) > 0 && descending[0] {
		desc = descending[0]
	}

	// sort descending
	canonicalized.SortVersions(out, desc)

	return out
}

func SortedVersions(versions []string, descending ...bool) []string {
	desc := false
	if len(descending) > 0 && descending[0] {
		desc = descending[0]
	}

	sortedParsedVersions := SortedParsedVersions(versions, desc)

	// flatten all parsed versions
	sortedVersions := make([]string, 0)
	for _, item := range sortedParsedVersions {
		sortedVersions = append(sortedVersions, item.Original)
	}

	return sortedVersions
}
