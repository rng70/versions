package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rng70/versions/vars"
)

/* ------------------------- */
/*      Ruby parser          */
/* ------------------------- */

// reRubyPart matches one Ruby Gemfile constraint token: op version
// The version may carry a SemVer pre-release suffix like "-preview.1.24081.5".
var reRubyPart = regexp.MustCompile(`^(~>|>=|<=|!=|>|<|=)?\s*([0-9]+(?:\.[0-9]+)*(?:-[A-Za-z0-9]+(?:\.[A-Za-z0-9]+)*)?)$`)

func ParseRuby(s string) ([][]vars.Constraint, error) {
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
		m := reRubyPart.FindStringSubmatch(part)
		if m == nil {
			continue
		}
		op := m[1]
		val := m[2]

		switch op {
		case "~>":
			// pessimistic constraint operator
			//   ~> 2.0   -> >= 2.0.0, < 3.0.0  (one dot: bump major)
			//   ~> 2.0.0 -> >= 2.0.0, < 2.1.0  (two dots: bump minor)
			nums := splitVersionNumsLegacy(val)
			lower := ensureThreePrerelease(val)
			var upper string
			if strings.Count(val, ".") <= 1 {
				upper = fmt.Sprintf("%d.0.0", nums[0]+1)
			} else {
				upper = fmt.Sprintf("%d.%d.0", nums[0], nums[1]+1)
			}
			ands = append(ands, vars.Constraint{Op: ">=", Ver: lower})
			ands = append(ands, vars.Constraint{Op: "<", Ver: ensureThree(upper)})
		case ">=", "<=", ">", "<", "!=":
			ands = append(ands, vars.Constraint{Op: op, Ver: ensureThreePrerelease(val)})
		default: // "=" or bare
			ands = append(ands, vars.Constraint{Op: "=", Ver: ensureThreePrerelease(val)})
		}
	}

	if len(ands) == 0 {
		return [][]vars.Constraint{}, nil
	}
	return [][]vars.Constraint{ands}, nil
}
