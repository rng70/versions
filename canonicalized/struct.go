package canonicalized

type TypeTag struct {
	Name string `json:"name"`
	Tag  int64  `json:"tag"`
}

type BuildMetadata struct {
	Tag string `json:"tag"`
}

type TimestampInfo struct {
	Original string `json:"original"`
	Parsed   string `json:"parsed"`
}

type CommitHashInfo struct {
	Original        string `json:"original"`
	Parsed          string `json:"parsed"`
	In              string `json:"in"`
	CommitsSinceTag *int64 `json:"commits_since_tag,omitempty"`
}

type Version struct {
	Prefix     string           `json:"prefix"`
	Major      *int64           `json:"major"`
	Minor      *int64           `json:"minor"`
	Patch      *int64           `json:"patch"`
	Revision   *int64           `json:"revision"`
	Type       []TypeTag        `json:"type"`
	Metadata   []BuildMetadata  `json:"build_metadata"`
	Extra      *int64           `json:"extra"`
	Canonical  string           `json:"canonical"`
	Original   string           `json:"original"`
	Timestamp  []TimestampInfo  `json:"timestamp"`
	CommitHash []CommitHashInfo `json:"commit_hash"`
}
