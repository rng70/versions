package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rng70/versions/vars"
)

/* ------------------------- */
/*      Rust/Cargo parser    */
/* ------------------------- */

// reRustPart matches one Cargo constraint token: op version
// The version may carry a SemVer pre-release suffix like "-preview.1.24081.5".
var reRustPart = regexp.MustCompile(`^(\^|~|>=|<=|!=|>|<|=)?\s*([0-9]+(?:\.[0-9*]+)*(?:-[A-Za-z0-9]+(?:\.[A-Za-z0-9]+)*)?)$`)

// caretRust implements Cargo caret semantics:
//
//	^1.2.3 -> >= 1.2.3, < 2.0.0
//	^0.2.3 -> >= 0.2.3, < 0.3.0
//	^0.0.3 -> >= 0.0.3, < 0.0.4
func caretRust(token string) []vars.Constraint {
	nums := splitVersionNumsLegacy(token)
	lower := ensureThreePrerelease(token)
	var upper string
	if nums[0] > 0 {
		upper = fmt.Sprintf("%d.0.0", nums[0]+1)
	} else if nums[1] > 0 {
		upper = fmt.Sprintf("0.%d.0", nums[1]+1)
	} else {
		upper = fmt.Sprintf("0.0.%d", nums[2]+1)
	}
	return []vars.Constraint{{Op: ">=", Ver: lower}, {Op: "<", Ver: upper}}
}

func ParseRust(s string) ([][]vars.Constraint, error) {
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

		// Wildcard forms: *, 1.*, 1.2.*
		if part == "*" || strings.HasSuffix(part, ".*") {
			ands = append(ands, expandWildcardNpm(part)...)
			continue
		}

		m := reRustPart.FindStringSubmatch(part)
		if m == nil {
			continue
		}
		op := m[1]
		val := m[2]

		switch op {
		case "^":
			ands = append(ands, caretRust(val)...)
		case "~":
			// ~1.2.3 -> >= 1.2.3, < 1.3.0
			// ~1.2   -> >= 1.2.0, < 1.3.0
			// ~1     -> >= 1.0.0, < 2.0.0
			nums := splitVersionNumsLegacy(val)
			lower := ensureThreePrerelease(val)
			var upper string
			if strings.Count(val, ".") == 0 {
				upper = fmt.Sprintf("%d.0.0", nums[0]+1)
			} else {
				upper = fmt.Sprintf("%d.%d.0", nums[0], nums[1]+1)
			}
			ands = append(ands, vars.Constraint{Op: ">=", Ver: lower})
			ands = append(ands, vars.Constraint{Op: "<", Ver: ensureThree(upper)})
		case ">=", "<=", ">", "<", "!=":
			ands = append(ands, vars.Constraint{Op: op, Ver: ensureThreePrerelease(val)})
		case "=":
			ands = append(ands, vars.Constraint{Op: "=", Ver: ensureThreePrerelease(val)})
		default: // bare version = exact constraint
			ands = append(ands, vars.Constraint{Op: "=", Ver: ensureThreePrerelease(val)})
		}
	}

	if len(ands) == 0 {
		return [][]vars.Constraint{}, nil
	}
	return [][]vars.Constraint{ands}, nil
}
