package semver

import "github.com/rng70/versions/canonicalized"

type Version struct {
	Original string                `json:"original_version_string"`
	Parsed   canonicalized.Version `json:"parsed_version_string"`
}

func NewVersionFromList(versions []string) []Version {
	var version []Version

	for _, s := range versions {
		p := canonicalized.ParseVersionString(s)

		version = append(version, Version{
			Original: s,
			Parsed:   p,
		})
	}

	return version
}

func NewVersion(s string) Version {
	p := canonicalized.ParseVersionString(s)

	return Version{
		Original: s,
		Parsed:   p,
	}
}
