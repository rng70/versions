package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/rng70/versions/v2/canonicalized"
	"github.com/rng70/versions/v2/vars"
)

/* ****************** Legacy utils ****************** */

func inc(version string, part string) string {
	nums := splitVersionNumsLegacy(version)
	switch part {
	case "major":
		return fmt.Sprintf("%d.0.0", nums[0]+1)
	case "minor":
		return fmt.Sprintf("%d.%d.0", nums[0], nums[1]+1)
	case "patch":
		return fmt.Sprintf("%d.%d.%d", nums[0], nums[1], nums[2]+1)
	default:
		return version
	}
}

func ensureThree(v string) string {
	nums, suffix := splitVersionNums(v)
	base := fmt.Sprintf("%d.%d.%d", nums[0], nums[1], nums[2])

	if suffix != "" {
		return base + suffix
	}
	return base
}

// ensureThreePrerelease is like ensureThree but preserves a full SemVer-style
// pre-release suffix (e.g. "-preview.1.24081.5").
// ensureThree alone truncates such suffixes because splitVersionNums stops at
// the first non-digit character inside a segment.
func ensureThreePrerelease(v string) string {
	if idx := strings.Index(v, "-"); idx >= 0 {
		return ensureThree(v[:idx]) + v[idx:]
	}
	return ensureThree(v)
}

// isBareVersion returns true when s is a bare version string with no leading
// operator — e.g. "9.0.0", "9.0", "9", or "9.0.0-preview.1.24081.5".
var reBareVersion = regexp.MustCompile(`^[0-9]+(?:\.[0-9]+)*(?:-[A-Za-z0-9]+(?:\.[A-Za-z0-9]+)*)?$`)

func isBareVersion(s string) bool { return reBareVersion.MatchString(s) }

// Expand "==1.2.*" into >=1.2.0 <1.3.0
func pyExpandWildcardEq(v string) []vars.Constraint {
	if strings.HasSuffix(v, ".*") {
		base := strings.TrimSuffix(v, ".*")
		nums := splitVersionNumsLegacy(base)
		lower := fmt.Sprintf("%d.%d.0", nums[0], nums[1])
		upper := fmt.Sprintf("%d.%d.0", nums[0], nums[1]+1)
		return []vars.Constraint{{Op: ">=", Ver: ensureThree(lower)}, {Op: "<", Ver: ensureThree(upper)}}
	}
	return []vars.Constraint{{Op: "=", Ver: ensureThreePrerelease(v)}}
}

// parsedConstraint holds a pre-parsed constraint to avoid repeated parsing.
type parsedConstraint struct {
	op  string
	ver canonicalized.Version
	raw string // original Ver string for "latest" check
}

// FilterMatches returns the subset of versions that satisfy at least one
// constraint group. Versions and constraints are parsed once upfront.
func FilterMatches(parsed [][]vars.Constraint, versions []string) []string {
	// Pre-parse all constraint versions once.
	groups := make([][]parsedConstraint, len(parsed))
	for i, ands := range parsed {
		pg := make([]parsedConstraint, len(ands))
		for j, c := range ands {
			pg[j] = parsedConstraint{op: c.Op, raw: c.Ver, ver: canonicalized.NewVersion(c.Ver)}
		}
		groups[i] = pg
	}

	// Pre-parse all candidate versions once.
	pvs := make([]canonicalized.Version, len(versions))
	for i, v := range versions {
		pvs[i] = canonicalized.NewVersion(v)
	}

	var out []string
	for i, v := range versions {
		for _, group := range groups {
			if satisfiesParsed(&pvs[i], v, group) {
				out = append(out, v)
				break
			}
		}
	}
	return out
}

// satisfiesParsed checks whether the pre-parsed version pv (original string v)
// satisfies every constraint in the AND group.
func satisfiesParsed(pv *canonicalized.Version, v string, ands []parsedConstraint) bool {
	for _, c := range ands {
		if c.raw == "latest" {
			return v == "latest"
		}
	}
	for _, c := range ands {
		if c.raw == "" {
			return false
		}
		cc := c.ver // local copy so we can take address
		switch c.op {
		case "=":
			if !pv.Equal(&cc) {
				return false
			}
		case "!=":
			if pv.Equal(&cc) {
				return false
			}
		case "<":
			if !pv.LessThan(&cc) {
				return false
			}
		case "<=":
			if !pv.LessThanOrEqual(&cc) {
				return false
			}
		case ">":
			if !pv.GreaterThan(&cc) {
				return false
			}
		case ">=":
			if !pv.GreaterThanOrEqual(&cc) {
				return false
			}
		case "<core":
			// Compare only major.minor.patch (ignore pre-release), used for compatible-release upper bound
			g := func(p *int64) int64 {
				if p == nil {
					return 0
				}
				return *p
			}
			vMaj, cMaj := g(pv.Major), g(cc.Major)
			vMin, cMin := g(pv.Minor), g(cc.Minor)
			vPat, cPat := g(pv.Patch), g(cc.Patch)
			coreNotLess := vMaj > cMaj ||
				(vMaj == cMaj && vMin > cMin) ||
				(vMaj == cMaj && vMin == cMin && vPat >= cPat)
			if coreNotLess {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func splitVersionNumsLegacy(v string) []int {
	// remove pre-release or build metadata, keep numeric prefix runs
	v = strings.SplitN(v, "-", 2)[0]
	v = strings.SplitN(v, "+", 2)[0]
	parts := strings.Split(v, ".")
	nums := make([]int, 0, 3)
	for _, p := range parts {
		if p == "" {
			nums = append(nums, 0)
			continue
		}

		i := 0
		for i < len(p) && p[i] >= '0' && p[i] <= '9' {
			i++
		}
		if i == 0 {
			nums = append(nums, 0)
			continue
		}
		n, _ := strconv.Atoi(p[:i])
		nums = append(nums, n)
	}
	for len(nums) < 3 {
		nums = append(nums, 0)
	}
	return nums
}

func splitVersionNums(v string) ([]int, string) {
	original := v

	// Strip build metadata (+...)
	if i := strings.Index(v, "+"); i != -1 {
		v = v[:i]
	}

	// Capture suffix after third numeric segment
	suffix := ""

	parts := strings.Split(v, ".")
	nums := make([]int, 0, 3)

	for idx, p := range parts {
		if len(nums) == 3 {
			// Everything after 3rd numeric part becomes suffix
			suffix = "." + strings.Join(parts[idx:], ".")
			break
		}

		if p == "" {
			nums = append(nums, 0)
			continue
		}

		i := 0
		for i < len(p) && p[i] >= '0' && p[i] <= '9' {
			i++
		}

		if i == 0 {
			nums = append(nums, 0)
			continue
		}

		n, _ := strconv.Atoi(p[:i])
		nums = append(nums, n)

		// Capture inline suffix like 3Final, including any remaining dot-separated parts
		if i < len(p) {
			rest := p[i:]
			remainParts := parts[idx+1:]
			if len(remainParts) > 0 {
				suffix = rest + "." + strings.Join(remainParts, ".")
			} else {
				suffix = rest
			}
			break
		}
	}

	// Pad missing numeric components
	for len(nums) < 3 {
		nums = append(nums, 0)
	}

	// Restore original suffix if nothing detected
	if suffix == "" && original != v {
		suffix = original[len(v):]
	}

	return nums, suffix
}

/* ****************** Legacy utils ****************** */

func SplitRequirement(req string) (string, string) {
	// Package name may include extras: abc[core]
	re := regexp.MustCompile(`^([a-zA-Z0-9._-]+(?:\[[a-zA-Z0-9._,-]+\])?)\s*(.*)$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(req))

	if len(matches) == 3 {
		return matches[1], strings.TrimSpace(matches[2])
	}
	return req, ""
}

func StringToInteger(s string) int {
	var n int
	_, _ = fmt.Sscanf(s, "%d", &n)
	return n
}
