package parser

import (
	"fmt"
	"strings"

	"github.com/rng70/versions/vars"
)

/* ------------------------- */
/*      Python parser        */
/* ------------------------- */

func ParsePython(s string) ([][]vars.Constraint, error) {
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
		m := vars.RePyPart.FindStringSubmatch(part)
		if m == nil {
			continue
		}
		op := m[1]
		val := m[2]
		switch op {
		case "==":
			ands = append(ands, pyExpandWildcardEq(val)...)
		case "===":
			ands = append(ands, vars.Constraint{Op: "=", Ver: ensureThreePrerelease(val)})
		case "!=":
			ands = append(ands, vars.Constraint{Op: "!=", Ver: ensureThreePrerelease(val)})
		case "<", "<=", ">", ">=":
			ands = append(ands, vars.Constraint{Op: op, Ver: ensureThreePrerelease(val)})
		case "~=":
			// compatible release operator
			// ~=1.4 -> >=1.4,<2.0
			// ~=1.4.5 -> >=1.4.5,<1.5.0
			nums := splitVersionNumsLegacy(val)
			lower := ensureThreePrerelease(val)
			var upper string
			if strings.Count(val, ".") == 1 {
				upper = fmt.Sprintf("%d.0.0", nums[0]+1)
			} else {
				upper = fmt.Sprintf("%d.%d.0", nums[0], nums[1]+1)
			}
			ands = append(ands, vars.Constraint{Op: ">=", Ver: lower})
			ands = append(ands, vars.Constraint{Op: "<core", Ver: ensureThree(upper)})
		default: // bare version (no operator) = exact match
			if val != "" {
				ands = append(ands, vars.Constraint{Op: "=", Ver: ensureThreePrerelease(val)})
			}
		}
	}
	if len(ands) == 0 {
		return [][]vars.Constraint{}, nil
	}
	return [][]vars.Constraint{ands}, nil
}
