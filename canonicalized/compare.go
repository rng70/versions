package canonicalized

func NewVersion(v string) Version {
	return newVersion(v)
}

func newVersion(v string) Version {
	return ParseVersionString(v)
}
