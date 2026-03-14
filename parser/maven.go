package parser

import (
	"regexp"
	"strings"

	"github.com/rng70/versions/vars"
)

/* ------------------------- */
/*      Maven parser         */
/* ------------------------- */

// reMavenExact matches a single exact version in brackets: [1.2.3] or [1.2.3-rc.1]
var reMavenExact = regexp.MustCompile(`^\s*\[\s*([0-9]+(?:\.[0-9]+)*(?:-[A-Za-z0-9]+(?:\.[A-Za-z0-9]+)*)?)\s*\]\s*$`)

func ParseMaven(s string) ([][]vars.Constraint, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return [][]vars.Constraint{}, nil
	}

	// Exact version: [1.2.3] or [1.2.3-rc.1]
	if m := reMavenExact.FindStringSubmatch(s); m != nil {
		return [][]vars.Constraint{{{Op: "=", Ver: ensureThreePrerelease(m[1])}}}, nil
	}

	// Try converting to bracket range.
	// This handles both bracket ranges ([lo,hi), etc.) and comparison-operator
	// forms like ">= 9.0.0-preview.1, <= 9.0.0-rc.1".
	rangeStr, err := ConstraintToRange(s)
	if err == nil && rangeStr != "" {
		if m := vars.ReNuGetRange.FindStringSubmatch(rangeStr); m != nil {
			open := m[1]
			lo := strings.TrimSpace(m[2])
			hi := strings.TrimSpace(m[3])
			close_ := m[4]
			var ands []vars.Constraint
			if lo != "" {
				if open == "[" {
					ands = append(ands, vars.Constraint{Op: ">=", Ver: ensureThreePrerelease(lo)})
				} else {
					ands = append(ands, vars.Constraint{Op: ">", Ver: ensureThreePrerelease(lo)})
				}
			}
			if hi != "" {
				if close_ == "]" {
					ands = append(ands, vars.Constraint{Op: "<=", Ver: ensureThreePrerelease(hi)})
				} else {
					ands = append(ands, vars.Constraint{Op: "<", Ver: ensureThreePrerelease(hi)})
				}
			}
			if len(ands) == 0 {
				return [][]vars.Constraint{}, nil
			}
			return [][]vars.Constraint{ands}, nil
		}
	}

	// Bare version (with optional pre-release) = exact constraint
	if isBareVersion(s) {
		return [][]vars.Constraint{{{Op: "=", Ver: ensureThreePrerelease(s)}}}, nil
	}

	return [][]vars.Constraint{}, nil
}
