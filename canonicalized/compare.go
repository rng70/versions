package canonicalized

func NewVersion(v string) Version {
	return newVersion(v)
}

func newVersion(v string) Version {
	return ParseVersionString(v)
}

// Compare compares this version to another version. This
// returns -1, 0, or 1 if this version is smaller, equal,
// or larger than the other version, respectively.
func (v *Version) Compare2(o *Version) int {
	// Check if the version is a exact match
	if v.Original == o.Original {
		return 0
	}

	// If the version string is not exactly same then compare
	// the major, minor, patch, revision and pre-release part

	// major release part
	if *v.Major > *o.Major {
		return 1
	}
	if *v.Major < *o.Major {
		return -1
	}

	// minor release part
	if *v.Minor > *o.Minor {
		return 1
	}
	if *v.Minor < *o.Minor {
		return -1
	}

	// patch
	if *v.Patch > *o.Patch {
		return 1
	}
	if *v.Patch < *o.Patch {
		return -1
	}

	// here the major, minor and patch portion is same
	// now only compare the prerelease part
	return v.compare(o)
}

func (v *Version) compare(o *Version) int {
	return 1
}
