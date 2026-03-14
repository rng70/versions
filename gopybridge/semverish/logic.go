package main

import (
	"encoding/json"

	"github.com/rng70/versions/v2/resolver"
	"github.com/rng70/versions/v2/semver"
	"github.com/rng70/versions/v2/vars"
)

// naturalSortedVersions is the pure-Go core of NaturalSortedVersions.
// Returns an error string prefixed with "ERROR: " on failure (matching the C export behaviour).
func naturalSortedVersions(input string, descending bool) string {
	var versions []string
	if err := json.Unmarshal([]byte(input), &versions); err != nil {
		return "ERROR: invalid json: " + err.Error()
	}
	res := semver.SortedVersions(versions, descending)
	out, _ := json.Marshal(res)
	return string(out)
}

// languageStyle maps a language/alias string to its vars.Style constant.
func languageStyle(language string) vars.Style {
	switch language {
	case "python", "py":
		return vars.StylePy
	case "nuget", "csharp", "dotnet", "cs":
		return vars.StyleNuGet
	case "npm", "node", "nodejs", "javascript", "js":
		return vars.StyleNPM
	case "maven", "java":
		return vars.StyleNPM
	default:
		return vars.StylePy
	}
}

// analyzeConstraints is the pure-Go core of AnalyzeConstrains.
// Returns an error string prefixed with "ERROR: " on failure.
func analyzeConstraints(language, constraints, versionsJSON string) string {
	var versions []string
	if err := json.Unmarshal([]byte(versionsJSON), &versions); err != nil {
		return "ERROR: invalid json: " + err.Error()
	}
	style := languageStyle(language)
	analysis := resolver.AnalyzeConstraint(style, constraints, versions)
	out, _ := json.Marshal(analysis.Matches)
	return string(out)
}
