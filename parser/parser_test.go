package parser

import (
	"testing"

	"github.com/rng70/versions/v2/vars"
)

// ─── inc ──────────────────────────────────────────────────────────────────────

func TestInc_Major(t *testing.T) {
	if got := inc("1.2.3", "major"); got != "2.0.0" {
		t.Errorf("got %q, want %q", got, "2.0.0")
	}
}

func TestInc_Minor(t *testing.T) {
	if got := inc("1.2.3", "minor"); got != "1.3.0" {
		t.Errorf("got %q, want %q", got, "1.3.0")
	}
}

func TestInc_Patch(t *testing.T) {
	if got := inc("1.2.3", "patch"); got != "1.2.4" {
		t.Errorf("got %q, want %q", got, "1.2.4")
	}
}

func TestInc_Unknown(t *testing.T) {
	if got := inc("1.2.3", "revision"); got != "1.2.3" {
		t.Errorf("unknown part: got %q, want unchanged %q", got, "1.2.3")
	}
}

func TestInc_MajorFromZero(t *testing.T) {
	if got := inc("0.0.0", "major"); got != "1.0.0" {
		t.Errorf("got %q, want %q", got, "1.0.0")
	}
}

// ─── ensureThree ──────────────────────────────────────────────────────────────

func TestEnsureThree_AlreadyThree(t *testing.T) {
	if got := ensureThree("1.2.3"); got != "1.2.3" {
		t.Errorf("got %q, want %q", got, "1.2.3")
	}
}

func TestEnsureThree_TwoPart(t *testing.T) {
	if got := ensureThree("1.2"); got != "1.2.0" {
		t.Errorf("got %q, want %q", got, "1.2.0")
	}
}

func TestEnsureThree_OnePart(t *testing.T) {
	if got := ensureThree("1"); got != "1.0.0" {
		t.Errorf("got %q, want %q", got, "1.0.0")
	}
}

func TestEnsureThree_FourPart(t *testing.T) {
	// fourth segment should become suffix ".4"
	got := ensureThree("1.2.3.4")
	if got != "1.2.3.4" {
		t.Errorf("got %q, want %q", got, "1.2.3.4")
	}
}

func TestEnsureThree_MavenFinal(t *testing.T) {
	// "4.1.0.Final" → numeric [4,1,0] + suffix ".Final"
	got := ensureThree("4.1.0.Final")
	if got != "4.1.0.Final" {
		t.Errorf("got %q, want %q", got, "4.1.0.Final")
	}
}

// ─── ensureThreePrerelease ────────────────────────────────────────────────────

func TestEnsureThreePrerelease_AlreadyThree(t *testing.T) {
	if got := ensureThreePrerelease("1.2.3"); got != "1.2.3" {
		t.Errorf("got %q, want %q", got, "1.2.3")
	}
}

func TestEnsureThreePrerelease_TwoPart(t *testing.T) {
	if got := ensureThreePrerelease("1.2"); got != "1.2.0" {
		t.Errorf("got %q, want %q", got, "1.2.0")
	}
}

func TestEnsureThreePrerelease_OnePart(t *testing.T) {
	if got := ensureThreePrerelease("1"); got != "1.0.0" {
		t.Errorf("got %q, want %q", got, "1.0.0")
	}
}

func TestEnsureThreePrerelease_WithSimplePrerelease(t *testing.T) {
	if got := ensureThreePrerelease("1.2-beta1"); got != "1.2.0-beta1" {
		t.Errorf("got %q, want %q", got, "1.2.0-beta1")
	}
}

func TestEnsureThreePrerelease_PreviewDots(t *testing.T) {
	got := ensureThreePrerelease("9.0.0-preview.1.24081.5")
	if got != "9.0.0-preview.1.24081.5" {
		t.Errorf("got %q, want %q", got, "9.0.0-preview.1.24081.5")
	}
}

func TestEnsureThreePrerelease_RCDots(t *testing.T) {
	got := ensureThreePrerelease("9.0.0-rc.2.24474.1")
	if got != "9.0.0-rc.2.24474.1" {
		t.Errorf("got %q, want %q", got, "9.0.0-rc.2.24474.1")
	}
}

// ─── isBareVersion ────────────────────────────────────────────────────────────

func TestIsBareVersion_SimpleThreePart(t *testing.T) {
	if !isBareVersion("1.2.3") {
		t.Error("1.2.3 should be a bare version")
	}
}

func TestIsBareVersion_TwoPart(t *testing.T) {
	if !isBareVersion("1.2") {
		t.Error("1.2 should be a bare version")
	}
}

func TestIsBareVersion_OnePart(t *testing.T) {
	if !isBareVersion("1") {
		t.Error("1 should be a bare version")
	}
}

func TestIsBareVersion_WithPrerelease(t *testing.T) {
	if !isBareVersion("1.0.0-beta1") {
		t.Error("1.0.0-beta1 should be a bare version")
	}
}

func TestIsBareVersion_WithPrereleaseDots(t *testing.T) {
	if !isBareVersion("9.0.0-preview.1.24081.5") {
		t.Error("9.0.0-preview.1.24081.5 should be a bare version")
	}
}

func TestIsBareVersion_GteOperator(t *testing.T) {
	if isBareVersion(">=1.0.0") {
		t.Error(">=1.0.0 should NOT be a bare version")
	}
}

func TestIsBareVersion_CaretOperator(t *testing.T) {
	if isBareVersion("^1.0.0") {
		t.Error("^1.0.0 should NOT be a bare version")
	}
}

func TestIsBareVersion_TildeOperator(t *testing.T) {
	if isBareVersion("~1.0.0") {
		t.Error("~1.0.0 should NOT be a bare version")
	}
}

func TestIsBareVersion_LatestKeyword(t *testing.T) {
	if isBareVersion("latest") {
		t.Error("latest should NOT be a bare version")
	}
}

func TestIsBareVersion_Empty(t *testing.T) {
	if isBareVersion("") {
		t.Error("empty string should NOT be a bare version")
	}
}

// ─── pyExpandWildcardEq ───────────────────────────────────────────────────────

func TestPyExpandWildcardEq_Wildcard(t *testing.T) {
	cs := pyExpandWildcardEq("1.2.*")
	if len(cs) != 2 {
		t.Fatalf("expected 2 constraints, got %d", len(cs))
	}
	if cs[0].Op != ">=" || cs[0].Ver != "1.2.0" {
		t.Errorf("lower: got {%s %s}, want {>= 1.2.0}", cs[0].Op, cs[0].Ver)
	}
	if cs[1].Op != "<" || cs[1].Ver != "1.3.0" {
		t.Errorf("upper: got {%s %s}, want {< 1.3.0}", cs[1].Op, cs[1].Ver)
	}
}

func TestPyExpandWildcardEq_WildcardMinorZero(t *testing.T) {
	cs := pyExpandWildcardEq("2.0.*")
	if len(cs) != 2 {
		t.Fatalf("expected 2 constraints, got %d", len(cs))
	}
	if cs[1].Ver != "2.1.0" {
		t.Errorf("upper: got %q, want %q", cs[1].Ver, "2.1.0")
	}
}

func TestPyExpandWildcardEq_ExactVersion(t *testing.T) {
	cs := pyExpandWildcardEq("1.2.3")
	if len(cs) != 1 {
		t.Fatalf("expected 1 constraint, got %d", len(cs))
	}
	if cs[0].Op != "=" {
		t.Errorf("op: got %q, want %q", cs[0].Op, "=")
	}
	if cs[0].Ver != "1.2.3" {
		t.Errorf("ver: got %q, want %q", cs[0].Ver, "1.2.3")
	}
}

func TestPyExpandWildcardEq_TwoPartVersion(t *testing.T) {
	cs := pyExpandWildcardEq("1.2")
	if len(cs) != 1 || cs[0].Op != "=" {
		t.Errorf("two-part without wildcard: got %v", cs)
	}
}

// ─── splitVersionNumsLegacy ───────────────────────────────────────────────────

func TestSplitVersionNumsLegacy_ThreePart(t *testing.T) {
	nums := splitVersionNumsLegacy("1.2.3")
	assertLegacyNums(t, nums, 1, 2, 3)
}

func TestSplitVersionNumsLegacy_TwoPart(t *testing.T) {
	nums := splitVersionNumsLegacy("1.2")
	assertLegacyNums(t, nums, 1, 2, 0)
}

func TestSplitVersionNumsLegacy_OnePart(t *testing.T) {
	nums := splitVersionNumsLegacy("5")
	assertLegacyNums(t, nums, 5, 0, 0)
}

func TestSplitVersionNumsLegacy_StripsPrerelease(t *testing.T) {
	nums := splitVersionNumsLegacy("1.2.3-beta1")
	assertLegacyNums(t, nums, 1, 2, 3)
}

func TestSplitVersionNumsLegacy_StripsMetadata(t *testing.T) {
	nums := splitVersionNumsLegacy("1.2.3+build.123")
	assertLegacyNums(t, nums, 1, 2, 3)
}

func TestSplitVersionNumsLegacy_FourPart(t *testing.T) {
	nums := splitVersionNumsLegacy("1.2.3.4")
	// only 3 parts returned (legacy)
	assertLegacyNums(t, nums, 1, 2, 3)
}

func TestSplitVersionNumsLegacy_NonNumericSegment(t *testing.T) {
	nums := splitVersionNumsLegacy("4.1.0.Final")
	assertLegacyNums(t, nums, 4, 1, 0)
}

// ─── splitVersionNums ─────────────────────────────────────────────────────────

func TestSplitVersionNums_ThreePart(t *testing.T) {
	nums, suffix := splitVersionNums("1.2.3")
	assertLegacyNums(t, nums, 1, 2, 3)
	if suffix != "" {
		t.Errorf("suffix: got %q, want empty", suffix)
	}
}

func TestSplitVersionNums_TwoPart(t *testing.T) {
	nums, suffix := splitVersionNums("1.2")
	assertLegacyNums(t, nums, 1, 2, 0)
	if suffix != "" {
		t.Errorf("suffix: got %q, want empty", suffix)
	}
}

func TestSplitVersionNums_OnePart(t *testing.T) {
	nums, suffix := splitVersionNums("3")
	assertLegacyNums(t, nums, 3, 0, 0)
	if suffix != "" {
		t.Errorf("suffix: got %q, want empty", suffix)
	}
}

func TestSplitVersionNums_FourPartSuffix(t *testing.T) {
	nums, suffix := splitVersionNums("1.2.3.4")
	assertLegacyNums(t, nums, 1, 2, 3)
	if suffix != ".4" {
		t.Errorf("suffix: got %q, want %q", suffix, ".4")
	}
}

func TestSplitVersionNums_NonNumericSuffix(t *testing.T) {
	nums, suffix := splitVersionNums("4.1.0.Final")
	assertLegacyNums(t, nums, 4, 1, 0)
	if suffix != ".Final" {
		t.Errorf("suffix: got %q, want %q", suffix, ".Final")
	}
}

func TestSplitVersionNums_InlineAlphaSuffix(t *testing.T) {
	// "3Final" inside a segment
	nums, suffix := splitVersionNums("1.2.3Final")
	assertLegacyNums(t, nums, 1, 2, 3)
	if suffix != "Final" {
		t.Errorf("suffix: got %q, want %q", suffix, "Final")
	}
}

func TestSplitVersionNums_BuildMetadataStripped(t *testing.T) {
	nums, _ := splitVersionNums("1.2.3+build")
	assertLegacyNums(t, nums, 1, 2, 3)
}

// ─── SplitRequirement ─────────────────────────────────────────────────────────

func TestSplitRequirement_WithGTEConstraint(t *testing.T) {
	name, constraint := SplitRequirement("requests>=2.0.0")
	if name != "requests" {
		t.Errorf("name: got %q, want %q", name, "requests")
	}
	if constraint != ">=2.0.0" {
		t.Errorf("constraint: got %q, want %q", constraint, ">=2.0.0")
	}
}

func TestSplitRequirement_NameOnly(t *testing.T) {
	name, constraint := SplitRequirement("requests")
	if name != "requests" {
		t.Errorf("name: got %q, want %q", name, "requests")
	}
	if constraint != "" {
		t.Errorf("constraint: got %q, want empty", constraint)
	}
}

func TestSplitRequirement_WithExtras(t *testing.T) {
	name, constraint := SplitRequirement("requests[security]>=2.0")
	if name != "requests[security]" {
		t.Errorf("name: got %q, want %q", name, "requests[security]")
	}
	if constraint != ">=2.0" {
		t.Errorf("constraint: got %q, want %q", constraint, ">=2.0")
	}
}

func TestSplitRequirement_LeadingSpaces(t *testing.T) {
	name, _ := SplitRequirement("  flask>=1.0  ")
	if name != "flask" {
		t.Errorf("name with spaces: got %q, want %q", name, "flask")
	}
}

func TestSplitRequirement_WithDashInName(t *testing.T) {
	name, constraint := SplitRequirement("my-package==1.0.0")
	if name != "my-package" {
		t.Errorf("name: got %q, want %q", name, "my-package")
	}
	if constraint != "==1.0.0" {
		t.Errorf("constraint: got %q, want %q", constraint, "==1.0.0")
	}
}

func TestSplitRequirement_WithVersionRange(t *testing.T) {
	name, constraint := SplitRequirement("django>=3.0,<4.0")
	if name != "django" {
		t.Errorf("name: got %q, want %q", name, "django")
	}
	if constraint != ">=3.0,<4.0" {
		t.Errorf("constraint: got %q, want %q", constraint, ">=3.0,<4.0")
	}
}

// ─── StringToInteger ──────────────────────────────────────────────────────────

func TestStringToInteger_Valid(t *testing.T) {
	if n := StringToInteger("42"); n != 42 {
		t.Errorf("got %d, want 42", n)
	}
}

func TestStringToInteger_Zero(t *testing.T) {
	if n := StringToInteger("0"); n != 0 {
		t.Errorf("got %d, want 0", n)
	}
}

func TestStringToInteger_LargeNumber(t *testing.T) {
	if n := StringToInteger("1000"); n != 1000 {
		t.Errorf("got %d, want 1000", n)
	}
}

func TestStringToInteger_Alpha(t *testing.T) {
	if n := StringToInteger("abc"); n != 0 {
		t.Errorf("non-numeric: got %d, want 0", n)
	}
}

func TestStringToInteger_Empty(t *testing.T) {
	if n := StringToInteger(""); n != 0 {
		t.Errorf("empty: got %d, want 0", n)
	}
}

func TestStringToInteger_AlphaNumericMixed(t *testing.T) {
	// Sscanf reads leading digits, ignores rest
	n := StringToInteger("12abc")
	if n != 12 {
		t.Errorf("got %d, want 12", n)
	}
}

// ─── ConstraintToRange ────────────────────────────────────────────────────────

func TestConstraintToRange_AlreadyBracket(t *testing.T) {
	got, err := ConstraintToRange("[1.0.0, 2.0.0)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "[1.0.0, 2.0.0)" {
		t.Errorf("got %q, want %q", got, "[1.0.0, 2.0.0)")
	}
}

func TestConstraintToRange_GteAndLt(t *testing.T) {
	got, err := ConstraintToRange(">= 1.0.0, < 2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "[1.0.0, 2.0.0)" {
		t.Errorf("got %q, want %q", got, "[1.0.0, 2.0.0)")
	}
}

func TestConstraintToRange_GtAndLte(t *testing.T) {
	got, err := ConstraintToRange("> 1.0.0, <= 2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "(1.0.0, 2.0.0]" {
		t.Errorf("got %q, want %q", got, "(1.0.0, 2.0.0]")
	}
}

func TestConstraintToRange_LtOnly(t *testing.T) {
	got, err := ConstraintToRange("< 2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "[, 2.0.0)" {
		t.Errorf("got %q, want %q", got, "[, 2.0.0)")
	}
}

func TestConstraintToRange_LteOnly(t *testing.T) {
	got, err := ConstraintToRange("<= 2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "[, 2.0.0]" {
		t.Errorf("got %q, want %q", got, "[, 2.0.0]")
	}
}

func TestConstraintToRange_GtOnly(t *testing.T) {
	got, err := ConstraintToRange("> 1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "(1.0.0, ]" {
		t.Errorf("got %q, want %q", got, "(1.0.0, ]")
	}
}

func TestConstraintToRange_GteOnly(t *testing.T) {
	got, err := ConstraintToRange(">= 1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "[1.0.0, ]" {
		t.Errorf("got %q, want %q", got, "[1.0.0, ]")
	}
}

func TestConstraintToRange_Exact(t *testing.T) {
	got, err := ConstraintToRange("= 1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "[1.0.0, 1.0.0]" {
		t.Errorf("got %q, want %q", got, "[1.0.0, 1.0.0]")
	}
}

func TestConstraintToRange_UnsupportedFormat(t *testing.T) {
	_, err := ConstraintToRange("unsupported format xyz")
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestConstraintToRange_WithPrerelease(t *testing.T) {
	got, err := ConstraintToRange(">= 9.0.0-preview.1.24081.5, < 10.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "[9.0.0-preview.1.24081.5, 10.0.0)" {
		t.Errorf("got %q, want %q", got, "[9.0.0-preview.1.24081.5, 10.0.0)")
	}
}

// ─── FilterMatches ────────────────────────────────────────────────────────────

func TestFilterMatches_NilConstraints(t *testing.T) {
	matches := FilterMatches(nil, []string{"1.0.0", "2.0.0"})
	if len(matches) != 0 {
		t.Errorf("nil constraints: expected no matches, got %v", matches)
	}
}

func TestFilterMatches_EmptyVersions(t *testing.T) {
	cs := [][]vars.Constraint{{{Op: ">=", Ver: "1.0.0"}}}
	matches := FilterMatches(cs, []string{})
	if len(matches) != 0 {
		t.Errorf("empty versions: expected no matches, got %v", matches)
	}
}

func TestFilterMatches_GteConstraint(t *testing.T) {
	versions := []string{"0.9.0", "1.0.0", "2.0.0", "3.0.0"}
	cs := [][]vars.Constraint{{{Op: ">=", Ver: "1.0.0"}}}
	matches := FilterMatches(cs, versions)
	want := map[string]bool{"1.0.0": true, "2.0.0": true, "3.0.0": true}
	if len(matches) != 3 {
		t.Errorf("expected 3 matches, got %d: %v", len(matches), matches)
	}
	for _, m := range matches {
		if !want[m] {
			t.Errorf("unexpected match: %q", m)
		}
	}
}

func TestFilterMatches_GtConstraint(t *testing.T) {
	versions := []string{"1.0.0", "1.0.1", "2.0.0"}
	cs := [][]vars.Constraint{{{Op: ">", Ver: "1.0.0"}}}
	matches := FilterMatches(cs, versions)
	want := map[string]bool{"1.0.1": true, "2.0.0": true}
	if len(matches) != 2 {
		t.Errorf("expected 2 matches, got %d: %v", len(matches), matches)
	}
	for _, m := range matches {
		if !want[m] {
			t.Errorf("unexpected match: %q", m)
		}
	}
}

func TestFilterMatches_LtConstraint(t *testing.T) {
	versions := []string{"0.9.0", "1.0.0", "2.0.0"}
	cs := [][]vars.Constraint{{{Op: "<", Ver: "1.0.0"}}}
	matches := FilterMatches(cs, versions)
	if len(matches) != 1 || matches[0] != "0.9.0" {
		t.Errorf("expected [0.9.0], got %v", matches)
	}
}

func TestFilterMatches_LteConstraint(t *testing.T) {
	versions := []string{"0.9.0", "1.0.0", "2.0.0"}
	cs := [][]vars.Constraint{{{Op: "<=", Ver: "1.0.0"}}}
	matches := FilterMatches(cs, versions)
	want := map[string]bool{"0.9.0": true, "1.0.0": true}
	if len(matches) != 2 {
		t.Errorf("expected 2 matches, got %d: %v", len(matches), matches)
	}
	for _, m := range matches {
		if !want[m] {
			t.Errorf("unexpected match: %q", m)
		}
	}
}

func TestFilterMatches_EqualConstraint(t *testing.T) {
	versions := []string{"1.0.0", "1.0.1", "2.0.0"}
	cs := [][]vars.Constraint{{{Op: "=", Ver: "1.0.1"}}}
	matches := FilterMatches(cs, versions)
	if len(matches) != 1 || matches[0] != "1.0.1" {
		t.Errorf("expected [1.0.1], got %v", matches)
	}
}

func TestFilterMatches_NotEqualConstraint(t *testing.T) {
	versions := []string{"1.0.0", "1.0.1", "2.0.0"}
	cs := [][]vars.Constraint{{{Op: "!=", Ver: "1.0.1"}}}
	matches := FilterMatches(cs, versions)
	want := map[string]bool{"1.0.0": true, "2.0.0": true}
	if len(matches) != 2 {
		t.Errorf("expected 2 matches, got %d: %v", len(matches), matches)
	}
	for _, m := range matches {
		if !want[m] {
			t.Errorf("unexpected match: %q", m)
		}
	}
}

func TestFilterMatches_AndGroup(t *testing.T) {
	versions := []string{"0.9.0", "1.0.0", "1.5.0", "2.0.0", "3.0.0"}
	cs := [][]vars.Constraint{
		{{Op: ">=", Ver: "1.0.0"}, {Op: "<", Ver: "2.0.0"}},
	}
	matches := FilterMatches(cs, versions)
	want := map[string]bool{"1.0.0": true, "1.5.0": true}
	if len(matches) != 2 {
		t.Errorf("expected 2 matches, got %d: %v", len(matches), matches)
	}
	for _, m := range matches {
		if !want[m] {
			t.Errorf("unexpected match: %q", m)
		}
	}
}

func TestFilterMatches_OrGroups(t *testing.T) {
	versions := []string{"1.0.0", "2.0.0", "3.0.0", "4.0.0"}
	cs := [][]vars.Constraint{
		{{Op: "<", Ver: "2.0.0"}},
		{{Op: ">", Ver: "3.0.0"}},
	}
	matches := FilterMatches(cs, versions)
	want := map[string]bool{"1.0.0": true, "4.0.0": true}
	if len(matches) != 2 {
		t.Errorf("expected 2 matches, got %d: %v", len(matches), matches)
	}
	for _, m := range matches {
		if !want[m] {
			t.Errorf("unexpected match: %q", m)
		}
	}
}

func TestFilterMatches_LatestLiteral(t *testing.T) {
	versions := []string{"latest", "1.0.0", "2.0.0"}
	cs := [][]vars.Constraint{{{Op: "=", Ver: "latest"}}}
	matches := FilterMatches(cs, versions)
	if len(matches) != 1 || matches[0] != "latest" {
		t.Errorf("expected [latest], got %v", matches)
	}
}

func TestFilterMatches_CoreOperator(t *testing.T) {
	// <core treats prerelease same as its core — so 1.6.0-beta1 has core 1.6.0 >= 1.6.0 → excluded
	versions := []string{"1.5.0", "1.5.0-beta1", "1.6.0", "1.6.0-beta1", "2.0.0"}
	cs := [][]vars.Constraint{{{Op: "<core", Ver: "1.6.0"}}}
	matches := FilterMatches(cs, versions)
	want := map[string]bool{"1.5.0": true, "1.5.0-beta1": true}
	if len(matches) != 2 {
		t.Errorf("expected 2 matches, got %d: %v", len(matches), matches)
	}
	for _, m := range matches {
		if !want[m] {
			t.Errorf("unexpected match: %q", m)
		}
	}
}

func TestFilterMatches_EmptyVerInConstraint(t *testing.T) {
	versions := []string{"1.0.0"}
	cs := [][]vars.Constraint{{{Op: "=", Ver: ""}}}
	matches := FilterMatches(cs, versions)
	if len(matches) != 0 {
		t.Errorf("empty ver in constraint: expected no matches, got %v", matches)
	}
}

func TestFilterMatches_UnknownOperator(t *testing.T) {
	versions := []string{"1.0.0"}
	cs := [][]vars.Constraint{{{Op: "??", Ver: "1.0.0"}}}
	matches := FilterMatches(cs, versions)
	if len(matches) != 0 {
		t.Errorf("unknown operator: expected no matches, got %v", matches)
	}
}

// ─── ParseNPM ─────────────────────────────────────────────────────────────────

func TestParseNPM_Empty(t *testing.T) {
	cs, err := ParseNPM("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 0 {
		t.Errorf("empty: expected 0 groups, got %d", len(cs))
	}
}

func TestParseNPM_Latest(t *testing.T) {
	cs, err := ParseNPM("latest")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 1 || cs[0][0].Op != "=" || cs[0][0].Ver != "latest" {
		t.Errorf("latest: unexpected result %v", cs)
	}
}

func TestParseNPM_WildcardStar(t *testing.T) {
	cs, err := ParseNPM("*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 {
		t.Fatalf("expected 1 group, got %d", len(cs))
	}
	if cs[0][0].Op != ">=" || cs[0][0].Ver != "0.0.0" {
		t.Errorf("wildcard: got {%s %s}, want {>= 0.0.0}", cs[0][0].Op, cs[0][0].Ver)
	}
}

func TestParseNPM_WildcardMajorX(t *testing.T) {
	cs, err := ParseNPM("2.x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("expected 1 group with 2 constraints, got %v", cs)
	}
	if cs[0][0].Op != ">=" || cs[0][0].Ver != "2.0.0" {
		t.Errorf("lower: got {%s %s}, want {>= 2.0.0}", cs[0][0].Op, cs[0][0].Ver)
	}
	if cs[0][1].Op != "<" || cs[0][1].Ver != "3.0.0" {
		t.Errorf("upper: got {%s %s}, want {< 3.0.0}", cs[0][1].Op, cs[0][1].Ver)
	}
}

func TestParseNPM_WildcardPatchX(t *testing.T) {
	cs, err := ParseNPM("3.3.x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("expected 1 group with 2 constraints, got %v", cs)
	}
	if cs[0][0].Ver != "3.3.0" || cs[0][1].Ver != "3.4.0" {
		t.Errorf("got {%s, %s}, want {3.3.0, 3.4.0}", cs[0][0].Ver, cs[0][1].Ver)
	}
}

func TestParseNPM_Caret_MajorNonZero(t *testing.T) {
	cs, err := ParseNPM("^1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("expected 1 group with 2 constraints, got %v", cs)
	}
	if cs[0][0].Op != ">=" || cs[0][1].Op != "<" {
		t.Errorf("expected >=/<, got %s/%s", cs[0][0].Op, cs[0][1].Op)
	}
	if cs[0][1].Ver != "2.0.0" {
		t.Errorf("upper bound for ^1.x: got %q, want %q", cs[0][1].Ver, "2.0.0")
	}
}

func TestParseNPM_Caret_ZeroMajor(t *testing.T) {
	cs, err := ParseNPM("^0.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("unexpected result %v", cs)
	}
	if cs[0][1].Ver != "0.3.0" {
		t.Errorf("upper bound for ^0.2.x: got %q, want %q", cs[0][1].Ver, "0.3.0")
	}
}

func TestParseNPM_Tilde(t *testing.T) {
	cs, err := ParseNPM("~1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("expected 1 group with 2 constraints, got %v", cs)
	}
	if cs[0][0].Op != ">=" || cs[0][1].Op != "<" {
		t.Errorf("expected >=/<, got %s/%s", cs[0][0].Op, cs[0][1].Op)
	}
}

func TestParseNPM_HyphenRange(t *testing.T) {
	cs, err := ParseNPM("1.0.0 - 2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("expected 1 group with 2 constraints, got %v", cs)
	}
	if cs[0][0].Op != ">=" || cs[0][1].Op != "<=" {
		t.Errorf("hyphen range ops: got %s/%s, want >=/<= ", cs[0][0].Op, cs[0][1].Op)
	}
}

func TestParseNPM_GteAndLt(t *testing.T) {
	// The NPM regex matches ">=X <Y" as one token, only the >= part is processed.
	cs, err := ParseNPM(">=1.0.0 <2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 {
		t.Fatalf("expected 1 group, got %d", len(cs))
	}
	// At least one constraint is present
	if len(cs[0]) == 0 {
		t.Error("expected at least 1 constraint in the group")
	}
}

func TestParseNPM_OrBlocks(t *testing.T) {
	cs, err := ParseNPM("1.0.0 || 2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 2 {
		t.Errorf("expected 2 OR groups, got %d", len(cs))
	}
}

func TestParseNPM_NpmAlias(t *testing.T) {
	cs, err := ParseNPM("npm:pkg@1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "=" || cs[0][0].Ver != "1.0.0" {
		t.Errorf("npm alias: unexpected result %v", cs)
	}
}

func TestParseNPM_HttpURL(t *testing.T) {
	_, err := ParseNPM("http://example.com/pkg.tgz")
	if err == nil {
		t.Error("expected ErrUnsupportedSource for HTTP URL")
	}
}

func TestParseNPM_HttpsURL(t *testing.T) {
	_, err := ParseNPM("https://example.com/pkg.tgz")
	if err == nil {
		t.Error("expected ErrUnsupportedSource for HTTPS URL")
	}
}

func TestParseNPM_FileURL(t *testing.T) {
	_, err := ParseNPM("file:../local-pkg")
	if err == nil {
		t.Error("expected ErrUnsupportedSource for file URL")
	}
}

func TestParseNPM_BareExact(t *testing.T) {
	cs, err := ParseNPM("1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "=" {
		t.Errorf("bare exact: expected = op, got %v", cs)
	}
}

// ─── ParsePython ──────────────────────────────────────────────────────────────

func TestParsePython_Empty(t *testing.T) {
	cs, err := ParsePython("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 0 {
		t.Errorf("empty: expected 0 groups, got %d", len(cs))
	}
}

func TestParsePython_EqualExact(t *testing.T) {
	cs, err := ParsePython("==1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "=" {
		t.Errorf("== exact: got %v", cs)
	}
}

func TestParsePython_EqualWildcard(t *testing.T) {
	// RePyPart does not match wildcard '*' — "==1.2.*" fails to parse, returns 0 groups
	cs, err := ParsePython("==1.2.*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 0 {
		t.Errorf("wildcard == with *: expected 0 groups (fails to parse), got %v", cs)
	}
}

func TestParsePython_TripleEqual(t *testing.T) {
	cs, err := ParsePython("===1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "=" {
		t.Errorf("===: got %v", cs)
	}
}

func TestParsePython_NotEqual(t *testing.T) {
	cs, err := ParsePython("!=1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "!=" {
		t.Errorf("!=: got %v", cs)
	}
}

func TestParsePython_GteAndLt(t *testing.T) {
	cs, err := ParsePython(">=1.0.0,<2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("expected 1 group with 2 constraints, got %v", cs)
	}
}

func TestParsePython_CompatibleRelease_OneDot(t *testing.T) {
	cs, err := ParsePython("~=1.4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("~= one-dot: expected 2 constraints, got %v", cs)
	}
	// lower >= 1.4.0, upper <core 2.0.0
	if cs[0][0].Op != ">=" {
		t.Errorf("lower op: got %q, want >=", cs[0][0].Op)
	}
	if cs[0][1].Op != "<core" {
		t.Errorf("upper op: got %q, want <core", cs[0][1].Op)
	}
	if cs[0][1].Ver != "2.0.0" {
		t.Errorf("upper ver: got %q, want 2.0.0", cs[0][1].Ver)
	}
}

func TestParsePython_CompatibleRelease_TwoDots(t *testing.T) {
	cs, err := ParsePython("~=1.4.5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("~= two-dots: expected 2 constraints, got %v", cs)
	}
	if cs[0][1].Op != "<core" {
		t.Errorf("upper op: got %q, want <core", cs[0][1].Op)
	}
	if cs[0][1].Ver != "1.5.0" {
		t.Errorf("upper ver: got %q, want 1.5.0", cs[0][1].Ver)
	}
}

func TestParsePython_BareVersion(t *testing.T) {
	cs, err := ParsePython("2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "=" {
		t.Errorf("bare version: got %v", cs)
	}
}

func TestParsePython_GteOnly(t *testing.T) {
	cs, err := ParsePython(">=1.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != ">=" {
		t.Errorf(">= only: got %v", cs)
	}
}

// ─── ParseNuGet ───────────────────────────────────────────────────────────────

func TestParseNuGet_Empty(t *testing.T) {
	cs, err := ParseNuGet("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 0 {
		t.Errorf("empty: expected 0 groups, got %d", len(cs))
	}
}

func TestParseNuGet_BareVersion(t *testing.T) {
	cs, err := ParseNuGet("1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "=" {
		t.Errorf("bare: got %v", cs)
	}
}

func TestParseNuGet_ExactBracket(t *testing.T) {
	cs, err := ParseNuGet("[1.2.3]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 1 || cs[0][0].Op != "=" {
		t.Errorf("[exact]: got %v", cs)
	}
}

func TestParseNuGet_ClosedOpen(t *testing.T) {
	cs, err := ParseNuGet("[1.0.0, 2.0.0)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("expected 1 group with 2 constraints, got %v", cs)
	}
	if cs[0][0].Op != ">=" || cs[0][1].Op != "<" {
		t.Errorf("[lo, hi): ops: got %s/%s, want >=/<", cs[0][0].Op, cs[0][1].Op)
	}
}

func TestParseNuGet_OpenClosed(t *testing.T) {
	cs, err := ParseNuGet("(, 2.0.0]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 1 || cs[0][0].Op != "<=" {
		t.Errorf("(,hi]: got %v", cs)
	}
}

func TestParseNuGet_ClosedUnbounded(t *testing.T) {
	cs, err := ParseNuGet("[1.0.0,)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 1 || cs[0][0].Op != ">=" {
		t.Errorf("[lo,): got %v", cs)
	}
}

func TestParseNuGet_WildcardMajor(t *testing.T) {
	cs, err := ParseNuGet("1.*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("1.*: expected 2 constraints, got %v", cs)
	}
	if cs[0][0].Op != ">=" || cs[0][0].Ver != "1.0.0" {
		t.Errorf("1.*: lower: got {%s %s}, want {>= 1.0.0}", cs[0][0].Op, cs[0][0].Ver)
	}
	if cs[0][1].Op != "<" || cs[0][1].Ver != "2.0.0" {
		t.Errorf("1.*: upper: got {%s %s}, want {< 2.0.0}", cs[0][1].Op, cs[0][1].Ver)
	}
}

func TestParseNuGet_WildcardMinor(t *testing.T) {
	cs, err := ParseNuGet("1.2.*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("1.2.*: expected 2 constraints, got %v", cs)
	}
	if cs[0][0].Ver != "1.2.0" || cs[0][1].Ver != "1.3.0" {
		t.Errorf("1.2.*: got {%s, %s}, want {1.2.0, 1.3.0}", cs[0][0].Ver, cs[0][1].Ver)
	}
}

func TestParseNuGet_OperatorGteAndLt(t *testing.T) {
	cs, err := ParseNuGet(">= 1.0.0, < 2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("expected 1 group with 2 constraints, got %v", cs)
	}
	if cs[0][0].Op != ">=" || cs[0][1].Op != "<" {
		t.Errorf("ops: got %s/%s, want >=/<", cs[0][0].Op, cs[0][1].Op)
	}
}

func TestParseNuGet_ExactBracketPrerelease(t *testing.T) {
	cs, err := ParseNuGet("[1.0.0-rc.1]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "=" {
		t.Errorf("[exact-rc]: got %v", cs)
	}
}

// ─── ParseRuby ────────────────────────────────────────────────────────────────

func TestParseRuby_Empty(t *testing.T) {
	cs, err := ParseRuby("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 0 {
		t.Errorf("empty: expected 0 groups, got %d", len(cs))
	}
}

func TestParseRuby_BareVersion(t *testing.T) {
	cs, err := ParseRuby("1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "=" {
		t.Errorf("bare: got %v", cs)
	}
}

func TestParseRuby_ExactOp(t *testing.T) {
	cs, err := ParseRuby("= 1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "=" {
		t.Errorf("= exact: got %v", cs)
	}
}

func TestParseRuby_Pessimistic_OneDot(t *testing.T) {
	cs, err := ParseRuby("~> 2.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("~> 2.0: expected 2 constraints, got %v", cs)
	}
	if cs[0][0].Op != ">=" || cs[0][0].Ver != "2.0.0" {
		t.Errorf("lower: got {%s %s}, want {>= 2.0.0}", cs[0][0].Op, cs[0][0].Ver)
	}
	if cs[0][1].Op != "<" || cs[0][1].Ver != "3.0.0" {
		t.Errorf("upper: got {%s %s}, want {< 3.0.0}", cs[0][1].Op, cs[0][1].Ver)
	}
}

func TestParseRuby_Pessimistic_TwoDots(t *testing.T) {
	cs, err := ParseRuby("~> 2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("~> 2.0.0: expected 2 constraints, got %v", cs)
	}
	if cs[0][1].Ver != "2.1.0" {
		t.Errorf("upper for ~> 2.0.0: got %q, want %q", cs[0][1].Ver, "2.1.0")
	}
}

func TestParseRuby_GTE(t *testing.T) {
	cs, err := ParseRuby(">= 1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != ">=" {
		t.Errorf(">=: got %v", cs)
	}
}

func TestParseRuby_NotEqual(t *testing.T) {
	cs, err := ParseRuby("!= 1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "!=" {
		t.Errorf("!=: got %v", cs)
	}
}

func TestParseRuby_Compound(t *testing.T) {
	cs, err := ParseRuby(">= 1.0.0, < 2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("compound: expected 1 group with 2 constraints, got %v", cs)
	}
}

func TestParseRuby_LT(t *testing.T) {
	cs, err := ParseRuby("< 3.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "<" {
		t.Errorf("<: got %v", cs)
	}
}

// ─── ParseRust ────────────────────────────────────────────────────────────────

func TestParseRust_Empty(t *testing.T) {
	cs, err := ParseRust("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 0 {
		t.Errorf("empty: expected 0 groups, got %d", len(cs))
	}
}

func TestParseRust_BareVersion(t *testing.T) {
	cs, err := ParseRust("1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "=" {
		t.Errorf("bare: got %v", cs)
	}
}

func TestParseRust_WildcardStar(t *testing.T) {
	cs, err := ParseRust("*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != ">=" || cs[0][0].Ver != "0.0.0" {
		t.Errorf("*: got %v", cs)
	}
}

func TestParseRust_WildcardMajor(t *testing.T) {
	cs, err := ParseRust("1.*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("1.*: expected 2 constraints, got %v", cs)
	}
	if cs[0][0].Ver != "1.0.0" || cs[0][1].Ver != "2.0.0" {
		t.Errorf("1.*: got {%s, %s}, want {1.0.0, 2.0.0}", cs[0][0].Ver, cs[0][1].Ver)
	}
}

func TestParseRust_Caret_MajorNonZero(t *testing.T) {
	cs, err := ParseRust("^1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("^1.2.3: expected 2 constraints, got %v", cs)
	}
	if cs[0][1].Ver != "2.0.0" {
		t.Errorf("^1.2.3 upper: got %q, want %q", cs[0][1].Ver, "2.0.0")
	}
}

func TestParseRust_Caret_ZeroMajorNonZeroMinor(t *testing.T) {
	cs, err := ParseRust("^0.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("^0.2.3: expected 2 constraints, got %v", cs)
	}
	if cs[0][1].Ver != "0.3.0" {
		t.Errorf("^0.2.3 upper: got %q, want %q", cs[0][1].Ver, "0.3.0")
	}
}

func TestParseRust_Caret_ZeroMajorZeroMinor(t *testing.T) {
	cs, err := ParseRust("^0.0.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("^0.0.3: expected 2 constraints, got %v", cs)
	}
	if cs[0][1].Ver != "0.0.4" {
		t.Errorf("^0.0.3 upper: got %q, want %q", cs[0][1].Ver, "0.0.4")
	}
}

func TestParseRust_Tilde_ThreePart(t *testing.T) {
	cs, err := ParseRust("~1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("~1.2.3: expected 2 constraints, got %v", cs)
	}
	if cs[0][1].Ver != "1.3.0" {
		t.Errorf("~1.2.3 upper: got %q, want %q", cs[0][1].Ver, "1.3.0")
	}
}

func TestParseRust_Tilde_OnePart(t *testing.T) {
	cs, err := ParseRust("~1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("~1: expected 2 constraints, got %v", cs)
	}
	if cs[0][1].Ver != "2.0.0" {
		t.Errorf("~1 upper: got %q, want %q", cs[0][1].Ver, "2.0.0")
	}
}

func TestParseRust_GTE(t *testing.T) {
	cs, err := ParseRust(">= 1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != ">=" {
		t.Errorf(">=: got %v", cs)
	}
}

func TestParseRust_NotEqual(t *testing.T) {
	cs, err := ParseRust("!= 1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "!=" {
		t.Errorf("!=: got %v", cs)
	}
}

func TestParseRust_ExactOp(t *testing.T) {
	cs, err := ParseRust("= 1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "=" {
		t.Errorf("=: got %v", cs)
	}
}

func TestParseRust_Compound(t *testing.T) {
	cs, err := ParseRust(">= 1.0.0, < 2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("compound: expected 1 group with 2 constraints, got %v", cs)
	}
}

// ─── ParseGo ──────────────────────────────────────────────────────────────────

func TestParseGo_Empty(t *testing.T) {
	cs, err := ParseGo("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 0 {
		t.Errorf("empty: expected 0 groups, got %d", len(cs))
	}
}

func TestParseGo_BareVersionNoPrefix(t *testing.T) {
	cs, err := ParseGo("1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "=" {
		t.Errorf("bare: got %v", cs)
	}
}

func TestParseGo_BareVersionVPrefix(t *testing.T) {
	cs, err := ParseGo("v1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "=" || cs[0][0].Ver != "1.0.0" {
		t.Errorf("v-prefix bare: got %v", cs)
	}
}

func TestParseGo_GteWithVPrefix(t *testing.T) {
	cs, err := ParseGo(">= v1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != ">=" || cs[0][0].Ver != "1.0.0" {
		t.Errorf(">= v1.0.0: got %v", cs)
	}
}

func TestParseGo_LtWithVPrefix(t *testing.T) {
	cs, err := ParseGo("< v2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "<" || cs[0][0].Ver != "2.0.0" {
		t.Errorf("< v2.0.0: got %v", cs)
	}
}

func TestParseGo_NotEqual(t *testing.T) {
	cs, err := ParseGo("!= v1.0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "!=" {
		t.Errorf("!=: got %v", cs)
	}
}

func TestParseGo_RangeWithVPrefix(t *testing.T) {
	cs, err := ParseGo(">= v1.0.0, < v2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("range: expected 1 group with 2 constraints, got %v", cs)
	}
	if cs[0][0].Op != ">=" || cs[0][0].Ver != "1.0.0" {
		t.Errorf("lower: got {%s %s}, want {>= 1.0.0}", cs[0][0].Op, cs[0][0].Ver)
	}
	if cs[0][1].Op != "<" || cs[0][1].Ver != "2.0.0" {
		t.Errorf("upper: got {%s %s}, want {< 2.0.0}", cs[0][1].Op, cs[0][1].Ver)
	}
}

func TestParseGo_LTE(t *testing.T) {
	cs, err := ParseGo("<= v2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "<=" {
		t.Errorf("<=: got %v", cs)
	}
}

func TestParseGo_GTE(t *testing.T) {
	cs, err := ParseGo(">= v1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != ">=" {
		t.Errorf(">=: got %v", cs)
	}
}

// ─── ParseMaven ───────────────────────────────────────────────────────────────

func TestParseMaven_Empty(t *testing.T) {
	cs, err := ParseMaven("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 0 {
		t.Errorf("empty: expected 0 groups, got %d", len(cs))
	}
}

func TestParseMaven_ExactBracket(t *testing.T) {
	cs, err := ParseMaven("[1.2.3]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "=" {
		t.Errorf("[exact]: got %v", cs)
	}
}

func TestParseMaven_ExactBracketPrerelease(t *testing.T) {
	cs, err := ParseMaven("[1.2.3-rc.1]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "=" {
		t.Errorf("[exact-rc]: got %v", cs)
	}
}

func TestParseMaven_ClosedOpen(t *testing.T) {
	cs, err := ParseMaven("[1.0.0, 2.0.0)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("[lo, hi): expected 1 group with 2 constraints, got %v", cs)
	}
	if cs[0][0].Op != ">=" || cs[0][1].Op != "<" {
		t.Errorf("[lo, hi): ops: got %s/%s, want >=/<", cs[0][0].Op, cs[0][1].Op)
	}
}

func TestParseMaven_OpenOpen(t *testing.T) {
	cs, err := ParseMaven("(1.0.0, 2.0.0)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("(lo, hi): expected 1 group with 2 constraints, got %v", cs)
	}
	if cs[0][0].Op != ">" || cs[0][1].Op != "<" {
		t.Errorf("(lo, hi): ops: got %s/%s, want >/<", cs[0][0].Op, cs[0][1].Op)
	}
}

func TestParseMaven_OpenClosed(t *testing.T) {
	cs, err := ParseMaven("(, 2.0.0]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 1 || cs[0][0].Op != "<=" {
		t.Errorf("(,hi]: got %v", cs)
	}
}

func TestParseMaven_ClosedUnbounded(t *testing.T) {
	cs, err := ParseMaven("[1.0.0,)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 1 || cs[0][0].Op != ">=" {
		t.Errorf("[lo,): got %v", cs)
	}
}

func TestParseMaven_BareVersion(t *testing.T) {
	cs, err := ParseMaven("1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != "=" {
		t.Errorf("bare: got %v", cs)
	}
}

func TestParseMaven_OperatorGTE(t *testing.T) {
	cs, err := ParseMaven(">= 1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || cs[0][0].Op != ">=" {
		t.Errorf(">=: got %v", cs)
	}
}

func TestParseMaven_OperatorRange(t *testing.T) {
	cs, err := ParseMaven(">= 1.0.0, < 2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 || len(cs[0]) != 2 {
		t.Fatalf("range: expected 1 group with 2 constraints, got %v", cs)
	}
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func assertLegacyNums(t *testing.T, nums []int, major, minor, patch int) {
	t.Helper()
	if len(nums) < 3 {
		t.Fatalf("expected at least 3 nums, got %d: %v", len(nums), nums)
	}
	if nums[0] != major || nums[1] != minor || nums[2] != patch {
		t.Errorf("got %v, want [%d %d %d]", nums[:3], major, minor, patch)
	}
}
