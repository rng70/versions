package canonicalized

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func (v *Version) Compare(other *Version) int {
	// 1) core
	if d := cmpPtrInt(v.Major, other.Major); d != 0 {
		return d
	}
	if d := cmpPtrInt(v.Minor, other.Minor); d != 0 {
		return d
	}
	if d := cmpPtrInt(v.Patch, other.Patch); d != 0 {
		return d
	}
	if d := cmpPtrInt(v.Revision, other.Revision); d != 0 {
		return d
	}

	// 2) stable vs prerelease
	vPre := v.isPrerelease()
	oPre := other.isPrerelease()
	if vPre && !oPre {
		return -1
	}
	if !vPre {
		if oPre {
			return 1
		}

		return 0
	}

	// 3) compare prereleases
	return comparePrerelease(v, other)
}

func (v *Version) LessThan(other *Version) bool           { return v.Compare(other) < 0 }
func (v *Version) GreaterThan(other *Version) bool        { return v.Compare(other) > 0 }
func (v *Version) Equal(other *Version) bool              { return v.Compare(other) == 0 }
func (v *Version) LessThanOrEqual(other *Version) bool    { return v.Compare(other) <= 0 }
func (v *Version) GreaterThanOrEqual(other *Version) bool { return v.Compare(other) >= 0 }

func (v *Version) Prerelease() string {
	if !v.isPrerelease() {
		return ""
	}
	parts := make([]string, 0, len(v.Type)*2+1)
	for _, t := range v.Type {
		parts = append(parts, strings.ToLower(t.Name))
		if t.Tag != 0 {
			parts = append(parts, strconv.FormatInt(t.Tag, 10))
		}
	}
	if v.Extra != nil {
		parts = append(parts, strconv.FormatInt(*v.Extra, 10))
	}
	return strings.Join(parts, ".")
}

func (v *Version) MetadataStr() string {
	if len(v.Metadata) == 0 {
		return ""
	}
	tags := make([]string, 0, len(v.Metadata))
	for _, m := range v.Metadata {
		tags = append(tags, m.Tag)
	}
	return strings.Join(tags, ".")
}

func (v *Version) String() string {
	if strings.TrimSpace(v.Canonical) != "" {
		return v.Canonical
	}

	_maj := safeInt(v.Major)
	_min := safeInt(v.Minor)
	_pat := safeInt(v.Patch)

	base := fmt.Sprintf("%d.%d.%d", _maj, _min, _pat)
	if v.Revision != nil {
		base = base + "." + strconv.FormatInt(*v.Revision, 10)
	}
	if pre := v.Prerelease(); pre != "" {
		base = base + "-" + pre
	}
	if md := v.MetadataStr(); md != "" {
		base = base + "+" + md
	}
	if strings.TrimSpace(v.Prefix) != "" {
		return v.Prefix + base
	}
	return base
}

func (v *Version) CompareType(other *Version) bool {
	if len(v.Type) != len(other.Type) {
		return false
	}
	for i := range v.Type {
		if !strings.EqualFold(v.Type[i].Name, other.Type[i].Name) ||
			v.Type[i].Tag != other.Type[i].Tag {
			return false
		}
	}
	if (v.Extra == nil) != (other.Extra == nil) {
		return false
	}
	if v.Extra != nil && other.Extra != nil && *v.Extra != *other.Extra {
		return false
	}
	return true
}

// ===== Predicates =====

func (v *Version) IsStable() bool {
	if len(v.Type) == 0 {
		return true
	}
	for _, t := range v.Type {
		if strings.EqualFold(t.Name, "stable") {
			return true
		}
	}
	return false
}

func (v *Version) IsAlpha() bool {
	if v.IsStable() {
		return false
	}
	for _, t := range v.Type {
		if strings.EqualFold(t.Name, "alpha") {
			return true
		}
	}
	return false
}

func (v *Version) IsBeta() bool {
	if v.IsStable() {
		return false
	}
	for _, t := range v.Type {
		if strings.EqualFold(t.Name, "beta") {
			return true
		}
	}
	return false
}

func (v *Version) IsRC() bool {
	if v.IsStable() {
		return false
	}
	for _, t := range v.Type {
		if strings.EqualFold(t.Name, "rc") ||
			strings.EqualFold(t.Name, "release-candidate") {
			return true
		}
	}
	return false
}

func (v *Version) IsPreview() bool {
	if v.IsStable() {
		return false
	}
	for _, t := range v.Type {
		if strings.EqualFold(t.Name, "preview") {
			return true
		}
	}
	return false
}

func (v *Version) IsPseudo() bool {
	return len(v.Timestamp) > 0 && len(v.CommitHash) > 0
}

// ===== Internal Helpers =====

func (v *Version) isPrerelease() bool { return !v.IsStable() }

func safeInt(ptr *int64) int64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func cmpPtrInt(a, b *int64) int {
	// nil = 0
	av := safeInt(a)
	bv := safeInt(b)
	switch {
	case av < bv:
		return -1
	case av > bv:
		return 1
	default:
		return 0
	}
}

// ===== Pre-release comparison logic =====

func stageWeight(name string) int {
	switch strings.ToLower(name) {
	case "alpha", "a":
		return 10
	case "beta", "b":
		return 20
	case "preview", "pre", "prerelease":
		return 30
	case "rc", "release-candidate":
		return 40
	case "stable":
		return 100
	default:
		return 70
	}
}

type ident struct {
	num    *int64
	str    string
	weight int
}

func (id ident) isNum() bool { return id.num != nil }

func compareId(a, b ident) int {
	if a.isNum() && b.isNum() {
		return cmpPtrInt(a.num, b.num)
	}
	if a.isNum() && !b.isNum() {
		return -1
	}
	if !a.isNum() && b.isNum() {
		return 1
	}
	if d := cmpPtrInt(i64(int64(a.weight)), i64(int64(b.weight))); d != 0 {
		return d
	}
	as := strings.ToLower(a.str)
	bs := strings.ToLower(b.str)
	switch {
	case as < bs:
		return -1
	case as > bs:
		return 1
	default:
		return 0
	}
}

func comparePrerelease(v, o *Version) int {
	va := toPreIdents(v)
	vb := toPreIdents(o)

	n := len(va)
	if len(vb) < n {
		n = len(vb)
	}
	for i := 0; i < n; i++ {
		if d := compareId(va[i], vb[i]); d != 0 {
			return d
		}
	}
	switch {
	case len(va) < len(vb):
		return -1
	case len(va) > len(vb):
		return 1
	default:
		return 0
	}
}

func toPreIdents(v *Version) []ident {
	out := make([]ident, 0, len(v.Type)*2+1)
	for _, t := range v.Type {
		out = append(out, ident{str: strings.ToLower(t.Name), weight: stageWeight(t.Name)})
		if t.Tag != 0 {
			tag := t.Tag
			out = append(out, ident{num: &tag})
		}
	}
	if v.Extra != nil {
		out = append(out, ident{num: v.Extra})
	}
	return out
}

func i64(v int64) *int64 { return &v }

func SortVersions(versions []*Version, descending ...bool) {
	desc := false
	if len(descending) > 0 && descending[0] {
		desc = descending[0]
		fmt.Println("Found reverse", desc)
	}

	sort.Slice(versions, func(i, j int) bool {
		if desc {
			return versions[i].GreaterThan(versions[j])
		}
		return versions[i].LessThan(versions[j])
	})
}
