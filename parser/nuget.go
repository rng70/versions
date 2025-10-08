package parser

import (
	"fmt"
	"strings"

	"github.com/rng70/versions/vars"
)

/* ------------------------- */
/*     NuGet (C#) parser     */
/* ------------------------- */

func ParseNuGet(s string) [][]vars.Constraint {
	s = strings.TrimSpace(s)
	if s == "" {
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
		nums := splitVersionNums(base)
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
