package parser

import (
	"regexp"
	"strings"

	"github.com/rng70/versions/vars"
)

/* ------------------------- */
/*    Go modules parser      */
/* ------------------------- */

// reGoPart matches one Go module constraint token, stripping an optional v prefix.
// The version may carry a SemVer pre-release suffix like "-preview.1.24081.5".
var reGoPart = regexp.MustCompile(`^(>=|<=|!=|>|<|=)?\s*v?([0-9]+(?:\.[0-9]+)*(?:-[A-Za-z0-9]+(?:\.[A-Za-z0-9]+)*)?)$`)

func ParseGo(s string) ([][]vars.Constraint, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return [][]vars.Constraint{}, nil
	}

	parts := strings.Split(s, ",")
	var ands []vars.Constraint
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		m := reGoPart.FindStringSubmatch(part)
		if m == nil {
			continue
		}
		op := m[1]
		val := m[2] // v prefix already stripped by the regex

		switch op {
		case ">=", "<=", ">", "<", "!=":
			ands = append(ands, vars.Constraint{Op: op, Ver: ensureThreePrerelease(val)})
		default: // "=" or bare vX.Y.Z — treat as exact
			ands = append(ands, vars.Constraint{Op: "=", Ver: ensureThreePrerelease(val)})
		}
	}

	if len(ands) == 0 {
		return [][]vars.Constraint{}, nil
	}
	return [][]vars.Constraint{ands}, nil
}
