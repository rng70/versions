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

// SortedParsedVersions sorts a list of version strings and returns parsed Version structs.
//
// opts[0] = descending (default false)
// opts[1] = safeParse  (default true) — when true, versions with no numeric core
// (e.g. "reservation", "refs/heads/main") are excluded from the result.
func SortedParsedVersions(versions []string, opts ...bool) []*canonicalized.Version {
	out := NewVersionFromList(versions)

	desc := false
	safeParse := true
	if len(opts) > 0 {
		desc = opts[0]
	}
	if len(opts) > 1 {
		safeParse = opts[1]
	}

	if safeParse {
		filtered := out[:0]
		for _, v := range out {
			if v.Major != nil && *v.Major != -1 {
				filtered = append(filtered, v)
			}
		}
		out = filtered
	}

	canonicalized.SortVersions(out, desc)

	return out
}

// SortedVersions sorts a list of version strings and returns the original strings in sorted order.
//
// opts[0] = descending (default false)
// opts[1] = safeParse  (default true) — when true, versions with no numeric core
// (e.g. "reservation", "refs/heads/main") are excluded from the result.
func SortedVersions(versions []string, opts ...bool) []string {
	desc := false
	safeParse := true
	if len(opts) > 0 {
		desc = opts[0]
	}
	if len(opts) > 1 {
		safeParse = opts[1]
	}

	sortedParsedVersions := SortedParsedVersions(versions, desc, safeParse)

	sortedVersions := make([]string, 0, len(sortedParsedVersions))
	for _, item := range sortedParsedVersions {
		sortedVersions = append(sortedVersions, item.Original)
	}

	return sortedVersions
}
