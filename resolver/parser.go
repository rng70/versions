package resolver

import (
	"github.com/rng70/versions/parser"
	"github.com/rng70/versions/vars"
)

func AnalyzeConstraint(style vars.Style, constraint string, versions []string) vars.Analysis {
	raw := constraint
	switch style {
	case vars.StyleNPM:
		parsed, specialMatches, ok := parser.ParseNPM(constraint)
		if !ok && parsed == nil {
			// ignored source like http/file
			return vars.Analysis{Raw: raw, Parsed: nil, Matches: []string{}}
		}
		// If parseNPM returned a specialMatches (like "latest" or npm:pkg@v),
		// return that as matches (user requested those returned)
		if specialMatches != nil {
			return vars.Analysis{Raw: raw, Parsed: parsed, Matches: specialMatches}
		}
		// else filter against provided versions
		matches := parser.FilterMatches(parsed, versions)
		return vars.Analysis{Raw: raw, Parsed: parsed, Matches: matches}

	case vars.StylePy:
		parsed := parser.ParsePython(constraint)
		if len(parsed) == 0 {
			return vars.Analysis{Raw: raw, Parsed: parsed, Matches: []string{}}
		}
		matches := parser.FilterMatches(parsed, versions)
		return vars.Analysis{Raw: raw, Parsed: parsed, Matches: matches}

	case vars.StyleNuGet:
		parsed := parser.ParseNuGet(constraint)
		if len(parsed) == 0 {
			return vars.Analysis{Raw: raw, Parsed: parsed, Matches: []string{}}
		}
		matches := parser.FilterMatches(parsed, versions)
		return vars.Analysis{Raw: raw, Parsed: parsed, Matches: matches}

	default:
		return vars.Analysis{Raw: raw, Parsed: nil, Matches: []string{}}
	}
}
