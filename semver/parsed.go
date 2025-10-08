package semver

import "github.com/rng70/versions/canonicalized"

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
