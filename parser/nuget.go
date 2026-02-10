package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rng70/versions/vars"
)

/* ------------------------- */
/*     NuGet (C#) parser     */
/* ------------------------- */

var (
	rangePattern = regexp.MustCompile(
		`^(\[|\()\s*[^,]*\s*,\s*[^,\]]*\s*(\]|\))$`,
	)

	singleConstraintPattern = regexp.MustCompile(
		`^(<=|>=|<|>|=)\s*([0-9]+(?:\.[0-9]+)*)$`,
	)
)

type bound struct {
	version   string
	inclusive bool
}

func ConstraintToRange(input string) (string, error) {
	input = strings.TrimSpace(input)

	if rangePattern.MatchString(input) {
		return input, nil
	}

	parts := splitConstraints(input)

	var lower *bound
	var upper *bound

	for _, part := range parts {
		part = strings.TrimSpace(part)

		m := singleConstraintPattern.FindStringSubmatch(part)
		if m == nil {
			return "", fmt.Errorf("unsupported constraint format: %s", input)
		}

		op := m[1]
		version := m[2]

		switch op {
		case ">":
			lower = &bound{version, false}
		case ">=":
			lower = &bound{version, true}
		case "<":
			upper = &bound{version, false}
		case "<=":
			upper = &bound{version, true}
		case "=":
			lower = &bound{version, true}
			upper = &bound{version, true}
		}
	}

	return buildRange(lower, upper), nil
}

func splitConstraints(input string) []string {
	// supports: ">1.0,<2.0" and ">1.0, <2.0"
	return strings.Split(input, ",")
}

func buildRange(lower, upper *bound) string {
	leftBracket := "("
	rightBracket := ")"

	leftVersion := ""
	rightVersion := ""

	if lower != nil {
		leftVersion = lower.version
		if lower.inclusive {
			leftBracket = "["
		}
	} else {
		leftBracket = "["
	}

	if upper != nil {
		rightVersion = upper.version
		if upper.inclusive {
			rightBracket = "]"
		}
	} else {
		rightBracket = "]"
	}

	return fmt.Sprintf("%s%s, %s%s", leftBracket, leftVersion, rightVersion, rightBracket)
}

func ParseNuGet(s string) [][]vars.Constraint {
	s, err := ConstraintToRange(strings.TrimSpace(s))
	if s == "" || err != nil {
		return [][]vars.Constraint{}
	}
	// bracket range
	if m := vars.ReNuGetRange.FindStringSubmatch(s); m != nil {
		open := m[1]
		lo := strings.TrimSpace(m[2])
		hi := strings.TrimSpace(m[3])
		_close := m[4]
		var ands []vars.Constraint
		if lo != "" {
			if open == "[" {
				ands = append(ands, vars.Constraint{Op: ">=", Ver: ensureThree(lo)})
			} else {
				ands = append(ands, vars.Constraint{Op: ">", Ver: ensureThree(lo)})
			}
		}
		if hi != "" {
			if _close == "]" {
				ands = append(ands, vars.Constraint{Op: "<=", Ver: ensureThree(hi)})
			} else {
				ands = append(ands, vars.Constraint{Op: "<", Ver: ensureThree(hi)})
			}
		}
		return [][]vars.Constraint{ands}
	}
	// floating versions like 1.* or 1.2.*
	if strings.HasSuffix(s, ".*") {
		base := strings.TrimSuffix(s, ".*")
		nums := splitVersionNumsLegacy(base)
		if strings.Count(base, ".") == 0 {
			lower := fmt.Sprintf("%d.0.0", nums[0])
			upper := fmt.Sprintf("%d.0.0", nums[0]+1)
			return [][]vars.Constraint{{{Op: ">=", Ver: ensureThree(lower)}, {Op: "<", Ver: ensureThree(upper)}}}
		}
		lower := fmt.Sprintf("%d.%d.0", nums[0], nums[1])
		upper := fmt.Sprintf("%d.%d.0", nums[0], nums[1]+1)
		return [][]vars.Constraint{{{Op: ">=", Ver: ensureThree(lower)}, {Op: "<", Ver: ensureThree(upper)}}}
	}
	// bare numeric version => minimum
	if isNumericVersion(s) {
		return [][]vars.Constraint{{{Op: ">=", Ver: ensureThree(s)}}}
	}
	// unrecognized
	return [][]vars.Constraint{}
}
