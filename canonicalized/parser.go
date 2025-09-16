package canonicalized

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseVersionString(s string) Version {
	original := s
	s = removeNoise(s)

	start, end, core, coreHasDot := findCoreIndex(s)

	pv := Version{
		Prefix:     "",
		Type:       []TypeTag{},
		Metadata:   []BuildMetadata{},
		Canonical:  "",
		Original:   original,
		Timestamp:  []TimestampInfo{},
		CommitHash: []CommitHashInfo{},
	}

	if start == -1 {
		// No recognizable version core → treat whole string as prefix
		var one int64 = -1

		pv.Prefix = s
		pv.Major = &one
		pv.Minor = &one
		pv.Patch = &one
		pv.Revision = &one

		return pv
	}

	if start > 0 {
		pv.Prefix = s[:start]
	}
	after := s[end:]

	nums := parseCoreInts(core)
	setMajorMinorPatch(&pv, nums)

	// pre-release / build metadata
	pre := ""
	build := ""
	if i := strings.IndexByte(after, '+'); i >= 0 {
		pre = after[:i]
		build = after[i+1:]
	} else {
		pre = after
	}

	// sanitize pre: only start parsing after a separator (- . _)
	pre = strings.TrimLeft(pre, " ._")

	// parse the pre-tokens and extract all commits with tags or only commits
	if pre != "" {
		// extract all "-N-gHASH" occurrences in order from prePart
		pre = extractCommitsAndTags(pre, "pre", &pv)

		// then extract plain -gHASH occurrences
		pre = extractCommitsGOnly(pre, "pre", &pv)
	}

	// tokenize remaining prePart and
	// parse tokens into type/additional/timestamp/hex/numeric/extra
	if pre != "" {
		// tokenize by [-._]
		rawTokens := splitTokens(pre)

		var lastTypeIdx = -1
		var seenAnyType bool
		var setRevision = pv.Revision != nil

		for i, tok := range rawTokens {
			lower := tok
			// lower := strings.ToLower(tok)

			// try to parse the token to timestamp
			if reTimestamp14.MatchString(lower) || reTimestamp8.MatchString(lower) {
				pv.Timestamp = append(pv.Timestamp, TimestampInfo{
					Original: lower,
					Parsed:   parseTimestampToISO(lower),
				})
				continue
			}

			// pure hex commit (no leading g)
			if reHexCommit.MatchString(lower) {
				pv.CommitHash = append(pv.CommitHash, CommitHashInfo{
					Original:        lower,
					Parsed:          lower,
					CommitsSinceTag: nil,
				})
				continue
			}

			// if the token is numeric but revision part is not set
			if reNumeric.MatchString(lower) {
				// if first token after core AND no type seen, treat as revision (pre-release numeric like 0)
				if i == 0 && !seenAnyType && !setRevision {
					if n, err := strconv.ParseInt(lower, 10, 64); err == nil {
						pv.Revision = &n
						setRevision = true
						continue
					}
				}
				// if lastType exists and has tag==0, attach numeric as that type's tag
				if lastTypeIdx >= 0 && pv.Type[lastTypeIdx].Tag == 0 {
					if n, err := strconv.ParseInt(lower, 10, 64); err == nil {
						t := pv.Type[lastTypeIdx]
						t.Tag = n
						pv.Type[lastTypeIdx] = t
						continue
					}
				}
				// if numeric is last token, treat as extra
				if i == len(rawTokens)-1 && pv.Extra == nil {
					if n, err := strconv.ParseInt(lower, 10, 64); err == nil {
						pv.Extra = &n
						continue
					}
				}
				// fallback to additional
				pv.Metadata = append(pv.Metadata, BuildMetadata{Tag: lower})
				continue
			}

			// If token has trailing digits like alpha10 or rc1
			if a, d, ok := splitAlphaNumTail(lower); ok && a != "" && d != "" && isOnlyAlpha(a) {
				name := canonicalName(a)
				seenAnyType = true
				tagNum, _ := strconv.ParseInt(d, 10, 64)
				pv.Type = append(pv.Type, TypeTag{Name: name, Tag: tagNum})
				lastTypeIdx = len(pv.Type) - 1
				continue
			}

			// Pure alphabetic or mixed but alphabetic-leading
			alphaPrefix := takeAlphaPrefix(lower)
			if alphaPrefix != "" && isOnlyAlpha(alphaPrefix) {
				name := canonicalName(alphaPrefix)
				seenAnyType = true
				pv.Type = append(pv.Type, TypeTag{Name: name, Tag: 0})
				lastTypeIdx = len(pv.Type) - 1
				// If token had non-alpha rest (like "post20201221" already handled above),
				// keep the remainder as additional
				if rest := strings.TrimPrefix(lower, alphaPrefix); rest != "" {
					if isNumeric(rest) && lastTypeIdx >= 0 && pv.Type[lastTypeIdx].Tag == 0 {
						// Use numeric rest as tag for that type
						val, _ := strconv.ParseInt(rest, 10, 64)
						pv.Type[lastTypeIdx].Tag = val
					} else {
						pv.Metadata = append(pv.Metadata, BuildMetadata{Tag: rest})
					}
				}
				continue
			}

			// Fallback: keep it as additional (e.g., hashes like e7d0053e6)
			pv.Metadata = append(pv.Metadata, BuildMetadata{Tag: tok})
		}

		if !setRevision {
			var zero int64 = 0
			pv.Revision = &zero
		}
	}

	// build part: already preserved whole build as BuildMetadata with leading '+'
	// but also analyze build tokens for timestamps and hashes and attach them
	if build != "" {
		// extract all "-N-gHASH" occurrences in order from build
		build = extractCommitsAndTags(build, "build", &pv)

		// then extract plain -gHASH occurrences
		build = extractCommitsGOnly(build, "build", &pv)

		buildTokens := splitTokens(build)
		for _, bt := range buildTokens {
			if reTimestamp14.MatchString(bt) || reTimestamp8.MatchString(bt) {
				pv.Timestamp = append(pv.Timestamp, TimestampInfo{
					Original: bt,
					Parsed:   parseTimestampToISO(bt),
				})
				continue
			}
			if reHexCommit.MatchString(bt) {
				pv.CommitHash = append(pv.CommitHash, CommitHashInfo{
					Original:        bt,
					Parsed:          bt,
					CommitsSinceTag: nil,
				})
				continue
			}
			// else ignore for special parsing (it's already included as part of BuildMetadata)
			pv.Metadata = append(pv.Metadata, BuildMetadata{Tag: "+" + bt})
		}
	}

	// Build canonical semver-ish string
	_major, _minor, _patch := int64(0), int64(0), int64(0)
	if pv.Major != nil {
		_major = *pv.Major
	} else {
		// if we had no core at all (shouldn’t happen here), bail with empty canonical
		// but if we matched a pure int core (like "v3") where no dot was present,
		// we still use it as major
	}
	if pv.Minor != nil {
		_minor = *pv.Minor
	} else if coreHasDot {
		// core has dot but minor missing is rare; default to 0
		_minor = 0
	} else {
		// core does not have dot; default to 0
		_minor = 0
	}

	if pv.Patch != nil {
		_patch = *pv.Patch
	} else if coreHasDot {
		_patch = 0
	} else {
		// core has dot but patch missing is common; default to 0
		_patch = 0
	}

	canon := fmt.Sprintf("%d.%d.%d", _major, _minor, _patch)

	// Assemble prerelease tokens
	var preTokens []string
	if len(pv.Type) > 0 {
		for _, tt := range pv.Type {
			if tt.Tag > 0 {
				preTokens = append(preTokens, fmt.Sprintf("%s.%d", tt.Name, tt.Tag))
			} else {
				preTokens = append(preTokens, tt.Name)
			}
		}
	}
	if len(preTokens) == 0 && pv.Revision != nil {
		preTokens = append(preTokens, fmt.Sprintf("%d", *pv.Revision))
	}
	if pv.Extra != nil {
		preTokens = append(preTokens, fmt.Sprintf("%d", *pv.Extra))
	}
	if len(preTokens) > 0 {
		canon += "-" + strings.Join(preTokens, ".")
	}

	// Build metadata: take only additional tags that start with '+'
	var buildMeta []string
	for _, a := range pv.Metadata {
		if strings.HasPrefix(a.Tag, "+") {
			buildMeta = append(buildMeta, strings.TrimPrefix(a.Tag, "+"))
		}
	}
	if len(buildMeta) > 0 {
		canon += "+" + strings.Join(buildMeta, ".")
	}

	pv.Canonical = canon
	return pv
}
