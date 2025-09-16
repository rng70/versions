/********************************************************************************************************
*                                                                                                       *
* OPENREFACTORY CONFIDENTIAL                                                                            *
* --------------------------                                                                            *
*                                                                                                       *
* Copyright (c) 2025 OpenRefactory, Inc. All Rights Reserved.                                           *
*                                                                                                       *
* NOTICE: All information contained herein is, and remains the property of OpenRefactory, Inc. The      *
* intellectual and technical concepts contained herein are proprietary to OpenRefactory, Inc. and       *
* may be covered by U.S. and Foreign Patents, patents in process, and are protected by trade secret     *
* or copyright law. Dissemination of this information or reproduction of this material is strictly      *
* forbidden unless prior written permission is obtained from OpenRefactory, Inc.                        *
*                                                                                                       *
* Contributors: Al Arafat Tanin (OpenRefactory, Inc.)                                                   *
*********************************************************************************************************/

package resolver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/rng70/versions/canonicalized"
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

func filterMatches(parsed [][]Constraint, versions []string) []string {
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

func legacySatisfiesOne(v string, ands []Constraint) bool {
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

func satisfiesOne(v string, ands []Constraint) bool {
	// If any constraint is "latest", it only matches literal "latest"
	for _, c := range ands {
		if c.Ver == "latest" {
			return v == "latest"
		}
	}

	canonicalizedV := canonicalized.NewVersion(v)
	// Check numeric constraints
	for _, c := range ands {
		// if constraint value is not numeric-like, fail (except "latest" handled above)
		if c.Ver == "" {
			return false
		}

		canonicalizedC := canonicalized.NewVersion(c.Ver)
		switch c.Op {
		case "=":
			return canonicalizedV.Equal(&canonicalizedC)
		case "!=":
			return !canonicalizedV.Equal(&canonicalizedC)
		case "<":
			return canonicalizedV.LessThan(&canonicalizedC)
		case "<=":
			return canonicalizedV.LessThanOrEqual(&canonicalizedC)
		case ">":
			return canonicalizedV.GreaterThan(&canonicalizedC)
		case ">=":
			canonicalizedV.GreaterThanOrEqual(&canonicalizedC)
		default:
			return false
		}
	}
	return true
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
	fmt.Sscanf(s, "%d", &n)
	return n
}
