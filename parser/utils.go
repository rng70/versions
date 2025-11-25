package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/rng70/versions/canonicalized"
	"github.com/rng70/versions/vars"
)

/* ****************** Legacy utils ****************** */

func inc(version string, part string) string {
	nums := splitVersionNums(version)
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

func legacyCmpVersion(a, b string) int {
	aa := splitVersionNums(a)
	bb := splitVersionNums(b)
	n := len(aa)
	if len(bb) > n {
		n = len(bb)
	}
	for len(aa) < n {
		aa = append(aa, 0)
	}
	for len(bb) < n {
		bb = append(bb, 0)
	}
	for i := 0; i < n; i++ {
		if aa[i] < bb[i] {
			return -1
		}
		if aa[i] > bb[i] {
			return 1
		}
	}
	return 0
}

func ensureThree(v string) string {
	nums := splitVersionNums(v)
	return fmt.Sprintf("%d.%d.%d", nums[0], nums[1], nums[2])
}

func isNumericVersion(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !(r == '.' || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}

// Expand "==1.2.*" into >=1.2.0 <1.3.0
func pyExpandWildcardEq(v string) []vars.Constraint {
	if strings.HasSuffix(v, ".*") {
		base := strings.TrimSuffix(v, ".*")
		nums := splitVersionNums(base)
		lower := fmt.Sprintf("%d.%d.0", nums[0], nums[1])
		upper := fmt.Sprintf("%d.%d.0", nums[0], nums[1]+1)
		return []vars.Constraint{{Op: ">=", Ver: ensureThree(lower)}, {Op: "<", Ver: ensureThree(upper)}}
	}
	return []vars.Constraint{{Op: "=", Ver: ensureThree(v)}}
}

func FilterMatches(parsed [][]vars.Constraint, versions []string) []string {
	var out []string
	for _, v := range versions {
		for _, group := range parsed {
			if satisfiesOne(v, group) {
				out = append(out, v)
				break
			}
		}
	}
	return out
}

func legacySatisfiesOne(v string, ands []vars.Constraint) bool {
	// If any constraint is "latest", it only matches literal "latest"
	for _, c := range ands {
		if c.Ver == "latest" {
			return v == "latest"
		}
	}
	// Check numeric constraints
	for _, c := range ands {
		// if constraint value is not numeric-like, fail (except "latest" handled above)
		if c.Ver == "" {
			return false
		}
		switch c.Op {
		case "=":
			if legacyCmpVersion(v, c.Ver) != 0 {
				return false
			}
		case "!=":
			if legacyCmpVersion(v, c.Ver) == 0 {
				return false
			}
		case "<":
			if !(legacyCmpVersion(v, c.Ver) < 0) {
				return false
			}
		case "<=":
			if !(legacyCmpVersion(v, c.Ver) <= 0) {
				return false
			}
		case ">":
			if !(legacyCmpVersion(v, c.Ver) > 0) {
				return false
			}
		case ">=":
			if !(legacyCmpVersion(v, c.Ver) >= 0) {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func satisfiesOne(v string, ands []vars.Constraint) bool {
	// If any constraint is "latest", it only matches literal "latest"
	for _, c := range ands {
		if c.Ver == "latest" {
			return v == "latest"
		}
	}

	canonicalizedV := canonicalized.NewVersion(v)

	// Check numeric constraints
	result := true
	for _, c := range ands {
		// if constraint value is not numeric-like, fail (except "latest" handled above)
		if c.Ver == "" {
			return false
		}

		canonicalizedC := canonicalized.NewVersion(c.Ver)
		switch c.Op {
		case "=":
			result = result && canonicalizedV.Equal(&canonicalizedC)
		case "!=":
			result = result && !canonicalizedV.Equal(&canonicalizedC)
		case "<":
			result = result && canonicalizedV.LessThan(&canonicalizedC)
		case "<=":
			result = result && canonicalizedV.LessThanOrEqual(&canonicalizedC)
		case ">":
			result = result && canonicalizedV.GreaterThan(&canonicalizedC)
		case ">=":
			result = result && canonicalizedV.GreaterThanOrEqual(&canonicalizedC)
		default:
			result = result && false
		}
	}

	return result
}

func splitVersionNums(v string) []int {
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
