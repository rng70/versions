package vars

import "regexp"

type Constraint struct {
	Op  string
	Ver string
}

type Analysis struct {
	Raw     string
	Parsed  [][]Constraint
	Matches []string
}

type Style string

const (
	StyleNPM   Style = "npm"
	StylePy    Style = "python"
	StyleNuGet Style = "nuget"
)

type ConstraintResult struct {
	Raw     string
	Parsed  string
	Matches []string
}

// Notes: Go's regexp does not support non-capturing groups (?:...), so parentheses are used.
// We capture two main groups: operator (group 1) and token (group 2).
var (
	ReDashRange  = regexp.MustCompile(`^\s*([0-9]+(?:\.[0-9]+){2})\s*-\s*([0-9]+(?:\.[0-9]+){2})\s*$`)
	ReNpmToken   = regexp.MustCompile(`(?i)(<=|>=|<|>|=|~|\^)?\s*([0-9]+(\.[0-9]+){0,2}(\.(x|X|\*))?|latest|npm:[^\s@]+@\d+(\.\d+){0,2}|https?://\S+|file:[^\s]+|\*|[0-9]+)`)
	RePyPart     = regexp.MustCompile(`^(==|!=|<=|>=|<|>|~=|===)?\s*([0-9]+(\.[0-9]+){0,2}(\.\*)?)\s*$`)
	ReNuGetRange = regexp.MustCompile(`^\s*([\[\(])\s*([^,\s]*)\s*,\s*([^\]\)\s]*)\s*([\]\)])\s*$`)
)

var (
	TestVersions = []string{
		"1.0.0",
		"1.2.3",
		"1.2.4",
		"1.3.0",
		"2.0.1",
		"2.3.1",
		"2.4.4",
		"2.5.2",
		"2.9.9",
		"3.0.0",
		"latest",
	}
)
