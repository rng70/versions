package resolver

import (
	"errors"

	"github.com/rng70/versions/parser"
	"github.com/rng70/versions/vars"
)

func AnalyzeConstraint(style vars.Style, constraint string, versions []string) vars.Analysis {
	raw := constraint

	var parsed [][]vars.Constraint
	var err error

	switch style {
	case vars.StyleNPM:
		parsed, err = parser.ParseNPM(constraint)
	case vars.StylePy:
		parsed, err = parser.ParsePython(constraint)
	case vars.StyleNuGet:
		parsed, err = parser.ParseNuGet(constraint)
	case vars.StyleMaven:
		parsed, err = parser.ParseMaven(constraint)
	case vars.StyleRuby:
		parsed, err = parser.ParseRuby(constraint)
	case vars.StyleRust:
		parsed, err = parser.ParseRust(constraint)
	case vars.StyleGo:
		parsed, err = parser.ParseGo(constraint)
	default:
		return vars.Analysis{Raw: raw, Parsed: nil, Matches: []string{}}
	}

	if err != nil {
		if errors.Is(err, vars.ErrUnsupportedSource) {
			return vars.Analysis{Raw: raw, Parsed: nil, Matches: []string{}}
		}
		return vars.Analysis{Raw: raw, Parsed: nil, Matches: []string{}}
	}

	if len(parsed) == 0 {
		return vars.Analysis{Raw: raw, Parsed: parsed, Matches: []string{}}
	}

	matches := parser.FilterMatches(parsed, versions)
	return vars.Analysis{Raw: raw, Parsed: parsed, Matches: matches}
}
