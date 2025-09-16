















package resolver

import (
	"fmt"
	"strconv"
	"strings"
)

/* ------------------------- */
/*         NPM parser        */
/* ------------------------- */

func expandWildcardNpm(ver string) []Constraint {
	v := strings.ToLower(ver)
	// "*" matches any version
	if v == "*" {
		return []Constraint{{Op: ">=", Ver: "0.0.0"}}
	}
	// "2.x" or "1.*" or "3.3.x"
	parts := strings.Split(v, ".")
	// normalize length
	for len(parts) < 3 {
		parts = append(parts, "0")
	}
	if strings.HasSuffix(v, ".x") || strings.HasSuffix(v, ".*") || strings.Contains(v, "x") || strings.Contains(v, "*") {
		// handle single segment like "2" is not wildcard; only "2.x" or "2.*"
		if len(parts) >= 1 && (parts[1] == "x" || parts[1] == "*" || parts[2] == "x" || parts[2] == "*") {
			// if second segment is x -> 2.x or 2.*  => >=2.0.0 <3.0.0
			if parts[1] == "x" || parts[1] == "*" {
				major, _ := strconv.Atoi(parts[0])
				lower := fmt.Sprintf("%d.0.0", major)
				upper := fmt.Sprintf("%d.0.0", major+1)
				return []Constraint{{Op: ">=", Ver: lower}, {Op: "<", Ver: upper}}
			}
			// else 3.3.x
			if parts[2] == "x" || parts[2] == "*" {
				major, _ := strconv.Atoi(parts[0])
				minor, _ := strconv.Atoi(parts[1])
				lower := fmt.Sprintf("%d.%d.0", major, minor)
				upper := fmt.Sprintf("%d.%d.0", major, minor+1)
				return []Constraint{{Op: ">=", Ver: lower}, {Op: "<", Ver: upper}}
			}
		}
	}
	// fallback equality
	return []Constraint{{Op: "=", Ver: ver}}
}

func parseNPM(s string) ([][]Constraint, []string, bool) {
	s = strings.TrimSpace(s)

	// special literal: latest
	if s == "latest" {
		return [][]Constraint{{{Op: "=", Ver: "latest"}}}, []string{"latest"}, true
	}

	// npm:pkg@1.0.0 -> extract version and return parsed + matches with that version
	if strings.HasPrefix(s, "npm:") {
		at := strings.LastIndex(s, "@")
		if at > -1 && at+1 < len(s) {
			ver := s[at+1:]
			// return parsed constraint and the version itself as the match (per user's request)
			return [][]Constraint{{{Op: "=", Ver: ver}}}, []string{ver}, true
		}
		// malformed -> no parse
		return [][]Constraint{}, []string{}, true
	}

	// ignore http/file sources
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") || strings.HasPrefix(s, "file:") {
		return nil, nil, false
	}

	// dash range
	if m := reDashRange.FindStringSubmatch(s); m != nil {
		l := ensureThree(m[1])
		h := ensureThree(m[2])
		return [][]Constraint{{{Op: ">=", Ver: l}, {Op: "<=", Ver: h}}}, nil, true
	}

	// OR blocks
	blocks := strings.Split(s, "||")
	var out [][]Constraint
	for _, b := range blocks {
		b = strings.TrimSpace(b)
		if b == "" {
			continue
		}
		matches := reNpmToken.FindAllStringSubmatch(b, -1)
		if matches == nil || len(matches) == 0 {
			continue
		}
		var ands []Constraint
		for _, m := range matches {
			// m[1] = operator (optional), m[2] = token
			op := strings.TrimSpace(m[1])
			token := strings.TrimSpace(m[2])
			if token == "" {
				continue
			}
			lowTok := strings.ToLower(token)
			if strings.HasPrefix(lowTok, "http://") || strings.HasPrefix(lowTok, "https://") || strings.HasPrefix(lowTok, "file:") {
				return nil, nil, false
			}
			// wildcard/star handling
			if lowTok == "latest" {
				ands = append(ands, Constraint{Op: "=", Ver: "latest"})
				continue
			}
			if strings.HasPrefix(lowTok, "npm:") {
				// extract after @ if present
				at := strings.LastIndex(token, "@")
				if at > -1 && at+1 < len(token) {
					v := token[at+1:]
					ands = append(ands, Constraint{Op: "=", Ver: v})
				}
				continue
			}
			// operator handling
			switch op {
			case "~":
				// ~1.2 or ~1.2.3 -> >=lower < next minor
				var lower string
				if strings.Count(token, ".") == 1 {
					lower = token + ".0"
				} else {
					lower = token
				}
				lower = ensureThree(lower)
				ands = append(ands, Constraint{Op: ">=", Ver: lower})
				ands = append(ands, Constraint{Op: "<", Ver: inc(lower, "minor")})
			case "^":
				// caret semantics:
				nums := splitVersionNums(token)
				lower := ensureThree(token)
				var upper string
				if nums[0] > 0 {
					upper = fmt.Sprintf("%d.0.0", nums[0]+1)
				} else if nums[1] > 0 {
					upper = fmt.Sprintf("0.%d.0", nums[1]+1)
				} else {
					upper = fmt.Sprintf("0.0.%d", nums[2]+1)
				}
				ands = append(ands, Constraint{Op: ">=", Ver: lower})
				ands = append(ands, Constraint{Op: "<", Ver: upper})
			default:
				// handle wildcard forms like 2.x, 3.3.x, 1.* or star token "*"
				if strings.Contains(strings.ToLower(token), "x") || strings.Contains(token, "*") || token == "*" {
					exp := expandWildcardNpm(token)
					ands = append(ands, exp...)
				} else {
					if op == "" {
						ands = append(ands, Constraint{Op: "=", Ver: ensureThree(token)})
					} else {
						ands = append(ands, Constraint{Op: op, Ver: ensureThree(token)})
					}
				}
			}
		} // end token loop
		if len(ands) > 0 {
			out = append(out, ands)
		}
	} // end blocks
	if len(out) == 0 {
		// could not parse into constraints, but it is valid input - return empty parse
		return [][]Constraint{}, nil, true
	}
	return out, nil, true
}

/* ------------------------- */
/*      Python parser        */
/* ------------------------- */

// Expand "==1.2.*" into >=1.2.0 <1.3.0
func pyExpandWildcardEq(v string) []Constraint {
	if strings.HasSuffix(v, ".*") {
		base := strings.TrimSuffix(v, ".*")
		nums := splitVersionNums(base)
		lower := fmt.Sprintf("%d.%d.0", nums[0], nums[1])
		upper := fmt.Sprintf("%d.%d.0", nums[0], nums[1]+1)
		return []Constraint{{Op: ">=", Ver: ensureThree(lower)}, {Op: "<", Ver: ensureThree(upper)}}
	}
	return []Constraint{{Op: "=", Ver: ensureThree(v)}}
}

func parsePython(s string) [][]Constraint {
	s = strings.TrimSpace(s)
	if s == "" {
		return [][]Constraint{}
	}
	parts := strings.Split(s, ",")
	var ands []Constraint
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		m := rePyPart.FindStringSubmatch(part)
		if m == nil {
			continue
		}
		op := m[1]
		val := m[2]
		switch op {
		case "==":
			ands = append(ands, pyExpandWildcardEq(val)...)
		case "===":
			ands = append(ands, Constraint{Op: "=", Ver: ensureThree(val)})
		case "!=":
			ands = append(ands, Constraint{Op: "!=", Ver: ensureThree(val)})
		case "<", "<=", ">", ">=":
			ands = append(ands, Constraint{Op: op, Ver: ensureThree(val)})
		case "~=":
			// compatible release operator
			// ~=1.4 -> >=1.4,<2.0
			// ~=1.4.5 -> >=1.4.5,<1.5.0
			nums := splitVersionNums(val)
			lower := ensureThree(val)
			var upper string
			if strings.Count(val, ".") == 1 {
				upper = fmt.Sprintf("%d.0.0", nums[0]+1)
			} else {
				upper = fmt.Sprintf("%d.%d.0", nums[0], nums[1]+1)
			}
			ands = append(ands, Constraint{Op: ">=", Ver: lower})
			ands = append(ands, Constraint{Op: "<", Ver: ensureThree(upper)})
		}
	}
	if len(ands) == 0 {
		return [][]Constraint{}
	}
	return [][]Constraint{ands}
}

/* ------------------------- */
/*     NuGet (C#) parser     */
/* ------------------------- */

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

func parseNuGet(s string) [][]Constraint {
	s = strings.TrimSpace(s)
	if s == "" {
		return [][]Constraint{}
	}
	// bracket range
	if m := reNuGetRange.FindStringSubmatch(s); m != nil {
		open := m[1]
		lo := strings.TrimSpace(m[2])
		hi := strings.TrimSpace(m[3])
		_close := m[4]
		var ands []Constraint
		if lo != "" {
			if open == "[" {
				ands = append(ands, Constraint{Op: ">=", Ver: ensureThree(lo)})
			} else {
				ands = append(ands, Constraint{Op: ">", Ver: ensureThree(lo)})
			}
		}
		if hi != "" {
			if _close == "]" {
				ands = append(ands, Constraint{Op: "<=", Ver: ensureThree(hi)})
			} else {
				ands = append(ands, Constraint{Op: "<", Ver: ensureThree(hi)})
			}
		}
		return [][]Constraint{ands}
	}
	// floating versions like 1.* or 1.2.*
	if strings.HasSuffix(s, ".*") {
		base := strings.TrimSuffix(s, ".*")
		nums := splitVersionNums(base)
		if strings.Count(base, ".") == 0 {
			lower := fmt.Sprintf("%d.0.0", nums[0])
			upper := fmt.Sprintf("%d.0.0", nums[0]+1)
			return [][]Constraint{{{Op: ">=", Ver: ensureThree(lower)}, {Op: "<", Ver: ensureThree(upper)}}}
		}
		lower := fmt.Sprintf("%d.%d.0", nums[0], nums[1])
		upper := fmt.Sprintf("%d.%d.0", nums[0], nums[1]+1)
		return [][]Constraint{{{Op: ">=", Ver: ensureThree(lower)}, {Op: "<", Ver: ensureThree(upper)}}}
	}
	// bare numeric version => minimum
	if isNumericVersion(s) {
		return [][]Constraint{{{Op: ">=", Ver: ensureThree(s)}}}
	}
	// unrecognized
	return [][]Constraint{}
}

func AnalyzeConstraint(style Style, constraint string, versions []string) Analysis {
	raw := constraint
	switch style {
	case StyleNPM:
		parsed, specialMatches, ok := parseNPM(constraint)
		if !ok && parsed == nil {
			// ignored source like http/file
			return Analysis{Raw: raw, Parsed: nil, Matches: []string{}}
		}
		// If parseNPM returned a specialMatches (like "latest" or npm:pkg@v),
		// return that as matches (user requested those returned)
		if specialMatches != nil {
			return Analysis{Raw: raw, Parsed: parsed, Matches: specialMatches}
		}
		// else filter against provided versions
		matches := filterMatches(parsed, versions)
		return Analysis{Raw: raw, Parsed: parsed, Matches: matches}

	case StylePy:
		parsed := parsePython(constraint)
		if len(parsed) == 0 {
			return Analysis{Raw: raw, Parsed: parsed, Matches: []string{}}
		}
		matches := filterMatches(parsed, versions)
		return Analysis{Raw: raw, Parsed: parsed, Matches: matches}

	case StyleNuGet:
		parsed := parseNuGet(constraint)
		if len(parsed) == 0 {
			return Analysis{Raw: raw, Parsed: parsed, Matches: []string{}}
		}
		matches := filterMatches(parsed, versions)
		return Analysis{Raw: raw, Parsed: parsed, Matches: matches}

	default:
		return Analysis{Raw: raw, Parsed: nil, Matches: []string{}}
	}
}
