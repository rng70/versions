package canonicalized

import (
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// semver-ish with at least one dot, with optional leading v
var (
	// find core like v?X.Y(.Z){0,2}
	reCoreWithDot = regexp.MustCompile(`(?i)\d+\.\d+(?:\.\d+){0,2}`)
	//reCoreWithDotWithV  = regexp.MustCompile(`(?i)v?\d+\.\d+(?:\.\d+){0,2}`)
	reCorePureInt = regexp.MustCompile(`(?i)^(?:v)?\d+$`)
	//reGitDescribeNum    = regexp.MustCompile(`-([0-9]+)-g([0-9a-f]{7,40})`)
	reGitDescribeNumAll = regexp.MustCompile(`-([0-9]+)-g([0-9a-f]{7,40})`)
	reGitDescribeGAll   = regexp.MustCompile(`-g([0-9a-f]{7,40})`)
	reTimestamp14       = regexp.MustCompile(`^\d{14}$`)
	reTimestamp8        = regexp.MustCompile(`^\d{8}$`)
	reHexCommit         = regexp.MustCompile(`^[0-9a-f]{7,40}$`)
	reNumeric           = regexp.MustCompile(`^\d+$`)
	//reTypeAlphaNum = regexp.MustCompile(`^([A-Za-z]+)(\d*)$`)
)

func trimWeirdQuotes(s string) string {
	// remove leading and/or trailing typical weird quotes/spaces
	return strings.Trim(s, " \t\r\n\"'")
}

func removeNoise(s string) string {
	s = trimWeirdQuotes(s)

	// remove common noise prefixes users sometimes embed
	s = strings.ReplaceAll(s, `\u003cbr\u003e`, "")
	s = strings.TrimSpace(s)

	return s
}

// extract up to 4 initial dotted ints from a string like "1.2.3.4"
func parseCoreInts(core string) (nums []int64) {
	parts := strings.Split(core, ".")
	for i := 0; i < len(parts) && i < 4; i++ {
		p := parts[i]
		// strip leading v on first token if any
		if i == 0 && (len(p) > 0 && (p[0] == 'v' || p[0] == 'V')) {
			p = p[1:]
		}
		if p == "" {
			break
		}
		if n, err := strconv.ParseInt(p, 10, 64); err == nil {
			nums = append(nums, n)
		} else {
			break
		}
	}
	return
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func splitAlphaNumTail(s string) (alpha string, digits string, ok bool) {
	// returns alpha prefix + trailing digit suffix if present (e.g., "alpha10" -> "alpha","10",true)
	if s == "" {
		return "", "", false
	}
	i := len(s) - 1
	for i >= 0 && unicode.IsDigit(rune(s[i])) {
		i--
	}
	if i < 0 || i == len(s)-1 {
		return s, "", false
	}
	alpha = s[:i+1]
	digits = s[i+1:]
	return alpha, digits, true
}

func canonicalName(raw string) string {
	r := strings.ToLower(raw)
	if v, ok := typeAlias[r]; ok {
		return v
	}
	return r
}

func findCoreIndex(s string) (start, end int, core string, hasDot bool) {
	loc := reCoreWithDot.FindStringIndex(s)
	if loc != nil {
		return loc[0], loc[1], s[loc[0]:loc[1]], true
	}

	// Fallback: accept lone v<digits> or <digits> as core only if the whole string (after trimming prefix) starts with it
	// Find first digit or 'v' followed by digits that is a token boundary
	for i := 0; i < len(s); i++ {
		r := s[i]
		if r == 'v' || r == 'V' || unicode.IsDigit(rune(r)) {
			rest := s[i:]
			tok := ""
			j := 0
			if len(rest) > 0 && (rest[0] == 'v' || rest[0] == 'V') {
				j++
			}
			for j < len(rest) && unicode.IsDigit(rune(rest[j])) {
				j++
			}
			if j > 0 {
				tok = rest[:j]
			}
			if tok != "" && reCorePureInt.MatchString(tok) {
				return i, i + len(tok), tok, false
			}
		}
	}
	return -1, -1, "", false
}

func takeAlphaPrefix(s string) string {
	i := 0
	for i < len(s) && unicode.IsLetter(rune(s[i])) {
		i++
	}
	return s[:i]
}

func isOnlyAlpha(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func parseTimestampToISO(ts string) string {
	if reTimestamp14.MatchString(ts) {
		if t, err := time.ParseInLocation("20060102150405", ts, time.UTC); err == nil {
			return t.UTC().Format(time.RFC3339)
		}
		return ts
	}
	if reTimestamp8.MatchString(ts) {
		if t, err := time.ParseInLocation("20060102", ts, time.UTC); err == nil {
			// date-only -> midnight UTC
			return t.UTC().Format(time.RFC3339)
		}
		return ts
	}
	return ts
}

// split tokens by - . _ space and drop empty tokens
func splitTokens(s string) []string {
	f := func(r rune) bool {
		return r == '-' || r == '.' || r == '_' || r == ' '
	}
	parts := strings.FieldsFunc(s, f)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// Extract all "-N-gHASH" occurrences in order from prePart
func extractCommitsAndTags(prePart, in string, parsed *Version) string {
	matches := reGitDescribeNumAll.FindAllStringSubmatch(prePart, -1)
	if len(matches) > 0 {
		for _, m := range matches {
			// m[1]=N, m[2]=hash
			if len(m) >= 3 {
				if cnt, err := strconv.ParseInt(m[1], 10, 64); err == nil {
					origHash := "g" + m[2]
					h := CommitHashInfo{
						Original:        origHash,
						Parsed:          strings.TrimPrefix(origHash, "g"),
						CommitsSinceTag: &cnt,
						In:              in,
					}
					parsed.CommitHash = append(parsed.CommitHash, h)
				}
				// remove the matched substring to avoid double-processing
				prePart = strings.Replace(prePart, "-"+m[1]+"-g"+m[2], "", 1)
			}
		}
	}

	return prePart
}

// Extract all "-gHASH" occurrences in order from prePart
func extractCommitsGOnly(prePart, in string, parsed *Version) string {
	matchesG := reGitDescribeGAll.FindAllStringSubmatch(prePart, -1)
	if len(matchesG) > 0 {
		for _, m := range matchesG {
			if len(m) >= 2 {
				origHash := "g" + m[1]
				h := CommitHashInfo{
					Original:        origHash,
					Parsed:          strings.TrimPrefix(origHash, "g"),
					CommitsSinceTag: nil,
					In:              in,
				}
				parsed.CommitHash = append(parsed.CommitHash, h)
				prePart = strings.Replace(prePart, "-g"+m[1], "", 1)
			}
		}
	}

	return prePart
}

func setMajorMinorPatch(pv *Version, nums []int64) {
	var zero int64 = 0

	if len(nums) > 3 {
		pv.Revision = &nums[3]
		pv.Patch = &nums[2]
		pv.Minor = &nums[1]
		pv.Major = &nums[0]
	} else if len(nums) > 2 {
		pv.Patch = &nums[2]
		pv.Minor = &nums[1]
		pv.Major = &nums[0]
	} else if len(nums) > 1 {
		pv.Patch = &zero
		pv.Minor = &nums[1]
		pv.Major = &nums[0]
	} else if len(nums) > 0 {
		pv.Patch = &zero
		pv.Minor = &zero
		pv.Major = &nums[0]
	} else {
		pv.Patch = &zero
		pv.Minor = &zero
		pv.Major = &zero
	}
}
