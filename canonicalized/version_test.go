package canonicalized

import (
	"fmt"
	"testing"
)

// ─── ParseVersionString ───────────────────────────────────────────────────────

func TestVersion_Lang3_0(t *testing.T) {
	v := ParseVersionString("COMMON_LANG_3_0")
	assertCore(t, v, 3, 0, 0)

	if v.Canonical != "3.0.0" || *v.Revision != 0 {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "3.0.0")
	}

	if v.Original != "COMMON_LANG_3_0" {
		t.Errorf("original: got %q, want %q", v.Original, "COMMON_LANG_3_0")
	}
}

func TestVersion_Lang3_0_1(t *testing.T) {
	v := ParseVersionString("COMMON_LANG_3_0_1")
	assertCore(t, v, 3, 0, 1)

	if v.Canonical != "3.0.1" || *v.Revision != 0 {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "3.0.1")
	}

	if v.Original != "COMMON_LANG_3_0_1" {
		t.Errorf("original: got %q, want %q", v.Original, "COMMON_LANG_3_0_1")
	}
}

func TestVersion_Lang3_0_1_0(t *testing.T) {
	v := ParseVersionString("COMMON_LANG_3_0_1_0")
	assertCore(t, v, 3, 0, 1)

	if v.Canonical != "3.0.1" || *v.Revision != 0 {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "3.0.1")
	}

	if v.Original != "COMMON_LANG_3_0_1_0" {
		t.Errorf("original: got %q, want %q", v.Original, "COMMON_LANG_3_0_1_0")
	}
}

func TestVersion_Lang3_0_1_2(t *testing.T) {
	v := ParseVersionString("COMMON_LANG_3_0_1_2")
	assertCore(t, v, 3, 0, 1)

	if v.Canonical != "3.0.1-2" || *v.Revision != 2 {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "3.0.1-2")
	}

	if v.Original != "COMMON_LANG_3_0_1_2" {
		t.Errorf("original: got %q, want %q", v.Original, "COMMON_LANG_3_0_1_2")
	}
}

func TestParseVersionString_ThreePart(t *testing.T) {
	v := ParseVersionString("1.2.3")
	assertCore(t, v, 1, 2, 3)
	if v.Canonical != "1.2.3" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "1.2.3")
	}
	if v.Original != "1.2.3" {
		t.Errorf("original: got %q, want %q", v.Original, "1.2.3")
	}
}

func TestParseVersionString_OriginalPreserved(t *testing.T) {
	v := ParseVersionString("1.2.3")
	if v.Original != "1.2.3" {
		t.Errorf("original: got %q, want %q", v.Original, "1.2.3")
	}
	if v.Canonical != "1.2.3" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "1.2.3")
	}
}

func TestParseVersionString_VPrefix(t *testing.T) {
	v := ParseVersionString("v2.0.1")
	assertCore(t, v, 2, 0, 1)
	if v.Prefix != "v" || *v.Revision != 0 {
		t.Errorf("prefix: got %q, want %q", v.Prefix, "v")
	}
	if v.Canonical != "2.0.1" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "2.0.1")
	}
	if v.Original != "v2.0.1" {
		t.Errorf("original: got %q, want %q", v.Original, "v2.0.1")
	}
}

func TestParseVersionString_TwoPart(t *testing.T) {
	v := ParseVersionString("4.7")
	assertCore(t, v, 4, 7, 0)
	fmt.Printf("%+v", v)
	if safeInt(v.Major) != 4 || safeInt(v.Minor) != 7 || *v.Revision != 0 {
		t.Errorf("got %d.%d, want 4.7", safeInt(v.Major), safeInt(v.Minor))
	}
	if v.Canonical != "4.7.0" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "4.7.0")
	}
	if v.Original != "4.7" {
		t.Errorf("original: got %q, want %q", v.Original, "4.7")
	}
}

func TestParseVersionString_OnePart(t *testing.T) {
	v := ParseVersionString("5")
	assertCore(t, v, 5, 0, 0)
	if safeInt(v.Major) != 5 && *v.Revision != 0 {
		t.Errorf("major: got %d, want 5", safeInt(v.Major))
	}
	if v.Canonical != "5.0.0" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "5.0.0")
	}
	if v.Original != "5" {
		t.Errorf("original: got %q, want %q", v.Original, "5")
	}
}

func TestParseVersionString_FourPart_Revision(t *testing.T) {
	v := ParseVersionString("1.2.3.4")
	assertCore(t, v, 1, 2, 3)
	if v.Revision == nil || *v.Revision != 4 {
		t.Errorf("revision: got %v, want 4", v.Revision)
	}
	if v.Canonical != "1.2.3-4" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "1.2.3-4")
	}
	if v.Original != "1.2.3.4" {
		t.Errorf("original: got %q, want %q", v.Original, "1.2.3.4")
	}
}

func TestParseVersionString_EmptyString(t *testing.T) {
	v := ParseVersionString("")
	// empty → no core; all fields set to -1
	if safeInt(v.Major) != -1 {
		t.Errorf("empty string: major should be -1, got %d", safeInt(v.Major))
	}
	if v.Canonical != "" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "")
	}
	if v.Original != "" {
		t.Errorf("original: got %q, want %q", v.Original, "")
	}
}

func TestParseVersionString_NoVersionCore(t *testing.T) {
	v := ParseVersionString("latest")
	if safeInt(v.Major) != -1 {
		t.Errorf("no-core: major should be -1, got %d", safeInt(v.Major))
	}
	if v.Prefix != "latest" {
		t.Errorf("no-core: prefix should be %q, got %q", "latest", v.Prefix)
	}
	if v.Canonical != "" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "")
	}
	if v.Original != "latest" {
		t.Errorf("original: got %q, want %q", v.Original, "latest")
	}
}

func TestParseVersionString_BetaWithNumericTag(t *testing.T) {
	v := ParseVersionString("1.0.0-beta1")
	assertCore(t, v, 1, 0, 0)
	if len(v.Type) == 0 || v.Type[0].Name != "beta" {
		t.Errorf("type[0].Name: got %v, want beta", v.Type)
	}
	if len(v.Type) > 0 && v.Type[0].Tag != 1 {
		t.Errorf("type[0].Tag: got %d, want 1", v.Type[0].Tag)
	}
	if v.Canonical != "1.0.0-beta.1" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "1.0.0-beta.1")
	}
	if v.Original != "1.0.0-beta1" {
		t.Errorf("original: got %q, want %q", v.Original, "1.0.0-beta1")
	}
}

func TestParseVersionString_AlphaWithNumericTag(t *testing.T) {
	v := ParseVersionString("2.0.0-alpha2")
	if len(v.Type) == 0 || v.Type[0].Name != "alpha" {
		t.Errorf("type[0].Name: got %v, want alpha", v.Type)
	}
	if len(v.Type) > 0 && v.Type[0].Tag != 2 {
		t.Errorf("type[0].Tag: got %d, want 2", v.Type[0].Tag)
	}
	if v.Canonical != "2.0.0-alpha.2" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "2.0.0-alpha.2")
	}
	if v.Original != "2.0.0-alpha2" {
		t.Errorf("original: got %q, want %q", v.Original, "2.0.0-alpha2")
	}
}

func TestParseVersionString_RCWithNumericTag(t *testing.T) {
	v := ParseVersionString("3.0.0-rc1")
	if len(v.Type) == 0 || v.Type[0].Name != "rc" {
		t.Errorf("type[0].Name: got %v, want rc", v.Type)
	}
	if v.Canonical != "3.0.0-rc.1" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "3.0.0-rc.1")
	}
	if v.Original != "3.0.0-rc1" {
		t.Errorf("original: got %q, want %q", v.Original, "3.0.0-rc1")
	}
}

func TestParseVersionString_PreviewDotNotation(t *testing.T) {
	v := ParseVersionString("4.0.0-preview.1")
	if len(v.Type) == 0 || v.Type[0].Name != "preview" {
		t.Errorf("type[0].Name: got %v, want preview", v.Type)
	}
	if len(v.Type) > 0 && v.Type[0].Tag != 1 {
		t.Errorf("type[0].Tag: got %d, want 1", v.Type[0].Tag)
	}
	if v.Canonical != "4.0.0-preview.1" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "4.0.0-preview.1")
	}
	if v.Original != "4.0.0-preview.1" {
		t.Errorf("original: got %q, want %q", v.Original, "4.0.0-preview.1")
	}
}

func TestParseVersionString_ComplexPreview(t *testing.T) {
	v := ParseVersionString("9.0.0-preview.1.24081.5")
	assertCore(t, v, 9, 0, 0)
	if len(v.Type) == 0 || v.Type[0].Name != "preview" {
		t.Errorf("type[0].Name: got %v, want preview", v.Type)
	}
	if v.Canonical != "9.0.0-preview.1.24081.5" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "9.0.0-preview.1.24081.5")
	}
	if v.Original != "9.0.0-preview.1.24081.5" {
		t.Errorf("original: got %q, want %q", v.Original, "9.0.0-preview.1.24081.5")
	}
}

func TestParseVersionString_ComplexRC(t *testing.T) {
	v := ParseVersionString("9.0.0-rc.2.24474.1.12234")
	assertCore(t, v, 9, 0, 0)

	if *v.Patch != 0 {
		t.Errorf("patch: got %v, want 0", *v.Patch)
	}
	if len(v.Type) == 0 || v.Type[0].Name != "rc" {
		t.Errorf("type[0].Name: got %v, want rc", v.Type)
	}
	if v.Canonical != "9.0.0-rc.2.24474.1.12234" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "9.0.0-rc.2.24474.1.12234")
	}
	if v.Original != "9.0.0-rc.2.24474.1.12234" {
		t.Errorf("original: got %q, want %q", v.Original, "9.0.0-rc.2.24474.1.12234")
	}
}

func TestParseVersionString_BuildMetadata(t *testing.T) {
	v := ParseVersionString("1.0.0+build.123")
	assertCore(t, v, 1, 0, 0)
	// metadata should be captured
	if v.Original != "1.0.0+build.123" {
		t.Errorf("original not preserved: got %q", v.Original)
	}
	if v.Canonical != "1.0.0+build.123" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "1.0.0+build.123")
	}
}

func TestParseVersionString_MavenFinalSuffix(t *testing.T) {
	v := ParseVersionString("4.1.0.Final")
	assertCore(t, v, 4, 1, 0)
	if v.Canonical != "4.1.0-final" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "4.1.0-final")
	}
	if v.Original != "4.1.0.Final" {
		t.Errorf("original: got %q, want %q", v.Original, "4.1.0.Final")
	}
}

func TestParseVersionString_MavenFinalSuffixWithNumber(t *testing.T) {
	v := ParseVersionString("4.1.0.Final.418")
	assertCore(t, v, 4, 1, 0)

	if *v.Patch != 0 {
		t.Errorf("patch: got %v, want 0", *v.Patch)
	}
	if v.Canonical != "4.1.0-final.418" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "4.1.0-final.418")
	}
	if v.Original != "4.1.0.Final.418" {
		t.Errorf("original: got %q, want %q", v.Original, "4.1.0.Final")
	}
}

func TestParseVersionString_StableHasNoType(t *testing.T) {
	v := ParseVersionString("2.0.0")
	if len(v.Type) != 0 {
		t.Errorf("stable version: want empty Type, got %v", v.Type)
	}
	if v.Canonical != "2.0.0" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "2.0.0")
	}
	if v.Original != "2.0.0" {
		t.Errorf("original: got %q, want %q", v.Original, "2.0.0")
	}
}

func TestParseVersionString_AliasA_IsAlpha(t *testing.T) {
	// "a" alias maps to "alpha"
	v := ParseVersionString("1.0.0-a1")
	if len(v.Type) == 0 || v.Type[0].Name != "alpha" {
		t.Errorf("alias 'a' should map to 'alpha', got %v", v.Type)
	}
	if v.Canonical != "1.0.0-alpha.1" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "1.0.0-alpha.1")
	}
	if v.Original != "1.0.0-a1" {
		t.Errorf("original: got %q, want %q", v.Original, "1.0.0-a1")
	}
}

func TestParseVersionString_AliasB_IsBeta(t *testing.T) {
	v := ParseVersionString("1.0.0-b2")
	if len(v.Type) == 0 || v.Type[0].Name != "beta" {
		t.Errorf("alias 'b' should map to 'beta', got %v", v.Type)
	}
	if v.Canonical != "1.0.0-beta.2" {
		t.Errorf("canonical: got %q, want %q", v.Canonical, "1.0.0-beta.2")
	}
	if v.Original != "1.0.0-b2" {
		t.Errorf("original: got %q, want %q", v.Original, "1.0.0-b2")
	}
}

// ─── NewVersion ───────────────────────────────────────────────────────────────

func TestNewVersion_Simple(t *testing.T) {
	v := NewVersion("1.2.3")
	assertCore(t, v, 1, 2, 3)
}

func TestNewVersion_VPrefix(t *testing.T) {
	v := NewVersion("v3.0.0")
	if v.Prefix != "v" {
		t.Errorf("prefix: got %q, want v", v.Prefix)
	}
	assertCore(t, v, 3, 0, 0)
}

func TestNewVersion_Prerelease(t *testing.T) {
	v := NewVersion("1.0.0-rc1")
	if !v.IsRC() {
		t.Error("expected IsRC() true")
	}
}

func TestNewVersion_NoCore(t *testing.T) {
	v := NewVersion("latest")
	if safeInt(v.Major) != -1 {
		t.Errorf("no-core: expected major=-1, got %d", safeInt(v.Major))
	}
}

func TestNewVersion_PreservesOriginal(t *testing.T) {
	v := NewVersion("v2.0.0-beta1")
	if v.Original != "v2.0.0-beta1" {
		t.Errorf("original: got %q, want %q", v.Original, "v2.0.0-beta1")
	}
}

// ─── Compare ──────────────────────────────────────────────────────────────────

func TestCompare_EqualVersions(t *testing.T) {
	a, b := NewVersion("1.2.3"), NewVersion("1.2.3")
	if a.Compare(&b) != 0 {
		t.Error("equal versions should compare to 0")
	}
}

func TestCompare_MajorLess(t *testing.T) {
	a, b := NewVersion("1.0.0"), NewVersion("2.0.0")
	if a.Compare(&b) >= 0 {
		t.Error("1.0.0 should be less than 2.0.0")
	}
}

func TestCompare_MajorGreater(t *testing.T) {
	a, b := NewVersion("3.0.0"), NewVersion("1.0.0")
	if a.Compare(&b) <= 0 {
		t.Error("3.0.0 should be greater than 1.0.0")
	}
}

func TestCompare_MinorLess(t *testing.T) {
	a, b := NewVersion("1.1.0"), NewVersion("1.2.0")
	if a.Compare(&b) >= 0 {
		t.Error("1.1.0 should be less than 1.2.0")
	}
}

func TestCompare_MinorGreater(t *testing.T) {
	a, b := NewVersion("1.5.0"), NewVersion("1.3.0")
	if a.Compare(&b) <= 0 {
		t.Error("1.5.0 should be greater than 1.3.0")
	}
}

func TestCompare_PatchLess(t *testing.T) {
	a, b := NewVersion("1.0.1"), NewVersion("1.0.9")
	if a.Compare(&b) >= 0 {
		t.Error("1.0.1 should be less than 1.0.9")
	}
}

func TestCompare_PatchGreater(t *testing.T) {
	a, b := NewVersion("1.0.5"), NewVersion("1.0.2")
	if a.Compare(&b) <= 0 {
		t.Error("1.0.5 should be greater than 1.0.2")
	}
}

func TestCompare_StableGreaterThanPrerelease(t *testing.T) {
	stable, pre := NewVersion("1.0.0"), NewVersion("1.0.0-beta1")
	if stable.Compare(&pre) <= 0 {
		t.Error("stable 1.0.0 should be greater than 1.0.0-beta1")
	}
}

func TestCompare_PrereleaseEqual(t *testing.T) {
	a, b := NewVersion("1.0.0-beta1"), NewVersion("1.0.0-beta1")
	if a.Compare(&b) != 0 {
		t.Error("identical pre-release versions should be equal")
	}
}

func TestCompare_AlphaBeforeBeta(t *testing.T) {
	alpha, beta := NewVersion("1.0.0-alpha1"), NewVersion("1.0.0-beta1")
	if alpha.Compare(&beta) >= 0 {
		t.Error("alpha should be less than beta")
	}
}

func TestCompare_BetaBeforeRC(t *testing.T) {
	beta, rc := NewVersion("1.0.0-beta1"), NewVersion("1.0.0-rc1")
	if beta.Compare(&rc) >= 0 {
		t.Error("beta should be less than rc")
	}
}

func TestCompare_PreviewBeforeRC(t *testing.T) {
	preview, rc := NewVersion("1.0.0-preview.1"), NewVersion("1.0.0-rc.1")
	if preview.Compare(&rc) >= 0 {
		t.Error("preview should be less than rc")
	}
}

func TestCompare_RCBeforeStable(t *testing.T) {
	rc, stable := NewVersion("1.0.0-rc1"), NewVersion("1.0.0")
	if rc.Compare(&stable) >= 0 {
		t.Error("rc should be less than stable")
	}
}

func TestCompare_TwoPreviewNumbers(t *testing.T) {
	p1, p2 := NewVersion("9.0.0-preview.1.24081.5"), NewVersion("9.0.0-preview.2.24128.5")
	if p1.Compare(&p2) >= 0 {
		t.Error("preview.1 should be less than preview.2")
	}
}

func TestCompare_TwoRCNumbers(t *testing.T) {
	r1, r2 := NewVersion("9.0.0-rc.1.24452.1"), NewVersion("9.0.0-rc.2.24474.1")
	if r1.Compare(&r2) >= 0 {
		t.Error("rc.1 should be less than rc.2")
	}
}

func TestCompare_DifferentMajorSamePrerelease(t *testing.T) {
	a, b := NewVersion("2.0.0-beta1"), NewVersion("1.0.0-beta1")
	if a.Compare(&b) <= 0 {
		t.Error("2.0.0-beta1 should be greater than 1.0.0-beta1")
	}
}

// ─── LessThan ─────────────────────────────────────────────────────────────────

func TestLessThan_Lower(t *testing.T) {
	a, b := NewVersion("1.0.0"), NewVersion("2.0.0")
	if !a.LessThan(&b) {
		t.Error("1.0.0 should be LessThan 2.0.0")
	}
}

func TestLessThan_Higher(t *testing.T) {
	a, b := NewVersion("3.0.0"), NewVersion("2.0.0")
	if a.LessThan(&b) {
		t.Error("3.0.0 should NOT be LessThan 2.0.0")
	}
}

func TestLessThan_Equal(t *testing.T) {
	a, b := NewVersion("1.0.0"), NewVersion("1.0.0")
	if a.LessThan(&b) {
		t.Error("equal version should NOT be LessThan")
	}
}

func TestLessThan_PrereleaseVsStable(t *testing.T) {
	pre, stable := NewVersion("1.0.0-beta1"), NewVersion("1.0.0")
	if !pre.LessThan(&stable) {
		t.Error("pre-release should be LessThan stable")
	}
}

// ─── GreaterThan ──────────────────────────────────────────────────────────────

func TestGreaterThan_Higher(t *testing.T) {
	a, b := NewVersion("2.0.0"), NewVersion("1.0.0")
	if !a.GreaterThan(&b) {
		t.Error("2.0.0 should be GreaterThan 1.0.0")
	}
}

func TestGreaterThan_Lower(t *testing.T) {
	a, b := NewVersion("1.0.0"), NewVersion("2.0.0")
	if a.GreaterThan(&b) {
		t.Error("1.0.0 should NOT be GreaterThan 2.0.0")
	}
}

func TestGreaterThan_Equal(t *testing.T) {
	a, b := NewVersion("1.0.0"), NewVersion("1.0.0")
	if a.GreaterThan(&b) {
		t.Error("equal version should NOT be GreaterThan")
	}
}

func TestGreaterThan_StableVsPrerelease(t *testing.T) {
	stable, pre := NewVersion("1.0.0"), NewVersion("1.0.0-rc1")
	if !stable.GreaterThan(&pre) {
		t.Error("stable should be GreaterThan pre-release of same core")
	}
}

// ─── Equal ────────────────────────────────────────────────────────────────────

func TestEqual_SameVersion(t *testing.T) {
	a, b := NewVersion("1.2.3"), NewVersion("1.2.3")
	if !a.Equal(&b) {
		t.Error("identical versions should be Equal")
	}
}

func TestEqual_DifferentVersions(t *testing.T) {
	a, b := NewVersion("1.2.3"), NewVersion("1.2.4")
	if a.Equal(&b) {
		t.Error("different versions should NOT be Equal")
	}
}

func TestEqual_SamePrerelease(t *testing.T) {
	a, b := NewVersion("1.0.0-beta1"), NewVersion("1.0.0-beta1")
	if !a.Equal(&b) {
		t.Error("same pre-release versions should be Equal")
	}
}

func TestEqual_DifferentPrerelease(t *testing.T) {
	a, b := NewVersion("1.0.0-beta1"), NewVersion("1.0.0-beta2")
	if a.Equal(&b) {
		t.Error("different pre-release tags should NOT be Equal")
	}
}

func TestEqual_StableVsPrerelease(t *testing.T) {
	stable, pre := NewVersion("1.0.0"), NewVersion("1.0.0-rc1")
	if stable.Equal(&pre) {
		t.Error("stable and pre-release should NOT be Equal")
	}
}

// ─── LessThanOrEqual ──────────────────────────────────────────────────────────

func TestLessThanOrEqual_Less(t *testing.T) {
	a, b := NewVersion("1.0.0"), NewVersion("2.0.0")
	if !a.LessThanOrEqual(&b) {
		t.Error("1.0.0 should be LessThanOrEqual 2.0.0")
	}
}

func TestLessThanOrEqual_Equal(t *testing.T) {
	a, b := NewVersion("1.0.0"), NewVersion("1.0.0")
	if !a.LessThanOrEqual(&b) {
		t.Error("1.0.0 should be LessThanOrEqual 1.0.0")
	}
}

func TestLessThanOrEqual_Greater(t *testing.T) {
	a, b := NewVersion("3.0.0"), NewVersion("2.0.0")
	if a.LessThanOrEqual(&b) {
		t.Error("3.0.0 should NOT be LessThanOrEqual 2.0.0")
	}
}

// ─── GreaterThanOrEqual ───────────────────────────────────────────────────────

func TestGreaterThanOrEqual_Greater(t *testing.T) {
	a, b := NewVersion("2.0.0"), NewVersion("1.0.0")
	if !a.GreaterThanOrEqual(&b) {
		t.Error("2.0.0 should be GreaterThanOrEqual 1.0.0")
	}
}

func TestGreaterThanOrEqual_Equal(t *testing.T) {
	a, b := NewVersion("1.0.0"), NewVersion("1.0.0")
	if !a.GreaterThanOrEqual(&b) {
		t.Error("1.0.0 should be GreaterThanOrEqual 1.0.0")
	}
}

func TestGreaterThanOrEqual_Less(t *testing.T) {
	a, b := NewVersion("1.0.0"), NewVersion("2.0.0")
	if a.GreaterThanOrEqual(&b) {
		t.Error("1.0.0 should NOT be GreaterThanOrEqual 2.0.0")
	}
}

// ─── Prerelease ───────────────────────────────────────────────────────────────

func TestPrerelease_StableVersionEmpty(t *testing.T) {
	v := NewVersion("1.0.0")
	if v.Prerelease() != "" {
		t.Errorf("stable version Prerelease() should be empty, got %q", v.Prerelease())
	}
}

func TestPrerelease_BetaVersion(t *testing.T) {
	v := NewVersion("1.0.0-beta1")
	p := v.Prerelease()
	if p == "" {
		t.Error("beta version Prerelease() should not be empty")
	}
}

func TestPrerelease_AlphaWithTag(t *testing.T) {
	v := NewVersion("1.0.0-alpha2")
	p := v.Prerelease()
	if p == "" {
		t.Error("alpha version Prerelease() should not be empty")
	}
}

func TestPrerelease_RCWithTag(t *testing.T) {
	v := NewVersion("1.0.0-rc3")
	p := v.Prerelease()
	if p == "" {
		t.Error("rc version Prerelease() should not be empty")
	}
}

func TestPrerelease_PreviewWithDot(t *testing.T) {
	v := NewVersion("1.0.0-preview.1")
	p := v.Prerelease()
	if p == "" {
		t.Error("preview version Prerelease() should not be empty")
	}
}

func TestPrerelease_WithExtra(t *testing.T) {
	// manually set Extra to test the Extra branch
	extra := int64(5)
	v := Version{
		Type:  []TypeTag{{Name: "preview", Tag: 1}},
		Extra: &extra,
	}
	p := v.Prerelease()
	if p == "" {
		t.Error("version with Extra should have non-empty Prerelease()")
	}
}

// ─── MetadataStr ──────────────────────────────────────────────────────────────

func TestMetadataStr_Empty(t *testing.T) {
	v := Version{Metadata: []BuildMetadata{}}
	if v.MetadataStr() != "" {
		t.Errorf("empty metadata should return empty string, got %q", v.MetadataStr())
	}
}

func TestMetadataStr_NoMetadataField(t *testing.T) {
	v := Version{}
	if v.MetadataStr() != "" {
		t.Errorf("nil metadata should return empty string, got %q", v.MetadataStr())
	}
}

func TestMetadataStr_Single(t *testing.T) {
	v := Version{Metadata: []BuildMetadata{{Tag: "build.123"}}}
	if v.MetadataStr() != "build.123" {
		t.Errorf("got %q, want %q", v.MetadataStr(), "build.123")
	}
}

func TestMetadataStr_Multiple(t *testing.T) {
	v := Version{Metadata: []BuildMetadata{{Tag: "sha"}, {Tag: "abc123"}}}
	if v.MetadataStr() != "sha.abc123" {
		t.Errorf("got %q, want %q", v.MetadataStr(), "sha.abc123")
	}
}

// ─── String ───────────────────────────────────────────────────────────────────

func TestString_WithCanonical(t *testing.T) {
	v := Version{Canonical: "1.2.3"}
	if v.String() != "1.2.3" {
		t.Errorf("String() with canonical: got %q, want %q", v.String(), "1.2.3")
	}
}

func TestString_BuildFromComponents(t *testing.T) {
	maj, min, pat := int64(2), int64(3), int64(4)
	v := Version{Major: &maj, Minor: &min, Patch: &pat, Type: []TypeTag{}}
	if v.String() != "2.3.4" {
		t.Errorf("String() from components: got %q, want %q", v.String(), "2.3.4")
	}
}

func TestString_WithPrefix(t *testing.T) {
	maj, min, pat := int64(1), int64(0), int64(0)
	v := Version{Prefix: "v", Major: &maj, Minor: &min, Patch: &pat, Type: []TypeTag{}}
	s := v.String()
	if s != "v1.0.0" {
		t.Errorf("String() with prefix: got %q, want %q", s, "v1.0.0")
	}
}

func TestString_WithRevision(t *testing.T) {
	maj, min, pat, rev := int64(1), int64(2), int64(3), int64(4)
	v := Version{Major: &maj, Minor: &min, Patch: &pat, Revision: &rev, Type: []TypeTag{}}
	s := v.String()
	if s != "1.2.3.4" {
		t.Errorf("String() with revision: got %q, want %q", s, "1.2.3.4")
	}
}

func TestString_WithPrerelease(t *testing.T) {
	maj, min, pat := int64(1), int64(0), int64(0)
	v := Version{
		Major: &maj, Minor: &min, Patch: &pat,
		Type: []TypeTag{{Name: "beta", Tag: 1}},
	}
	s := v.String()
	if s != "1.0.0-beta.1" {
		t.Errorf("String() with prerelease: got %q, want %q", s, "1.0.0-beta.1")
	}
}

func TestString_WithMetadata(t *testing.T) {
	maj, min, pat := int64(1), int64(0), int64(0)
	v := Version{
		Major:    &maj,
		Minor:    &min,
		Patch:    &pat,
		Type:     []TypeTag{},
		Metadata: []BuildMetadata{{Tag: "build"}},
	}
	s := v.String()
	if s != "1.0.0+build" {
		t.Errorf("String() with metadata: got %q, want %q", s, "1.0.0+build")
	}
}

func TestString_ParsedVersionUsesCanonical(t *testing.T) {
	v := NewVersion("1.2.3")
	// parsed versions always have Canonical set
	if v.String() != v.Canonical {
		t.Errorf("String() should equal Canonical for parsed versions")
	}
}

// ─── CompareType ──────────────────────────────────────────────────────────────

func TestCompareType_Equal(t *testing.T) {
	a, b := NewVersion("1.0.0-beta1"), NewVersion("1.0.0-beta1")
	if !a.CompareType(&b) {
		t.Error("same pre-release type should be equal")
	}
}

func TestCompareType_DifferentTypeName(t *testing.T) {
	a, b := NewVersion("1.0.0-alpha1"), NewVersion("1.0.0-beta1")
	if a.CompareType(&b) {
		t.Error("different type names should NOT be equal")
	}
}

func TestCompareType_DifferentTag(t *testing.T) {
	a, b := NewVersion("1.0.0-beta1"), NewVersion("1.0.0-beta2")
	if a.CompareType(&b) {
		t.Error("different type tags should NOT be equal")
	}
}

func TestCompareType_DifferentLength(t *testing.T) {
	a, b := NewVersion("1.0.0-alpha1"), NewVersion("1.0.0")
	if a.CompareType(&b) {
		t.Error("different type lengths should NOT be equal")
	}
}

func TestCompareType_BothStable(t *testing.T) {
	a, b := NewVersion("1.0.0"), NewVersion("2.0.0")
	if !a.CompareType(&b) {
		t.Error("two stable versions should have equal (empty) type")
	}
}

func TestCompareType_DifferentExtra(t *testing.T) {
	e1, e2 := int64(1), int64(2)
	a := Version{Type: []TypeTag{{Name: "preview", Tag: 1}}, Extra: &e1}
	b := Version{Type: []TypeTag{{Name: "preview", Tag: 1}}, Extra: &e2}
	if a.CompareType(&b) {
		t.Error("different Extra values should NOT be equal")
	}
}

func TestCompareType_NilExtraVsSet(t *testing.T) {
	e := int64(1)
	a := Version{Type: []TypeTag{{Name: "preview", Tag: 1}}}
	b := Version{Type: []TypeTag{{Name: "preview", Tag: 1}}, Extra: &e}
	if a.CompareType(&b) {
		t.Error("nil Extra vs set Extra should NOT be equal")
	}
}

// ─── IsStable ─────────────────────────────────────────────────────────────────

func TestIsStable_EmptyType(t *testing.T) {
	v := NewVersion("1.0.0")
	if !v.IsStable() {
		t.Error("version with no type should be stable")
	}
}

func TestIsStable_ExplicitStableTag(t *testing.T) {
	v := Version{Type: []TypeTag{{Name: "stable"}}}
	if !v.IsStable() {
		t.Error("version with 'stable' type should be stable")
	}
}

func TestIsStable_BetaIsNotStable(t *testing.T) {
	v := NewVersion("1.0.0-beta1")
	if v.IsStable() {
		t.Error("beta version should NOT be stable")
	}
}

func TestIsStable_AlphaIsNotStable(t *testing.T) {
	v := NewVersion("1.0.0-alpha1")
	if v.IsStable() {
		t.Error("alpha version should NOT be stable")
	}
}

func TestIsStable_RCIsNotStable(t *testing.T) {
	v := NewVersion("1.0.0-rc1")
	if v.IsStable() {
		t.Error("rc version should NOT be stable")
	}
}

func TestIsStable_PreviewIsNotStable(t *testing.T) {
	v := NewVersion("1.0.0-preview.1")
	if v.IsStable() {
		t.Error("preview version should NOT be stable")
	}
}

// ─── IsAlpha ──────────────────────────────────────────────────────────────────

func TestIsAlpha_Alpha(t *testing.T) {
	v := NewVersion("1.0.0-alpha1")
	if !v.IsAlpha() {
		t.Error("expected IsAlpha() true")
	}
}

func TestIsAlpha_StableNotAlpha(t *testing.T) {
	v := NewVersion("1.0.0")
	if v.IsAlpha() {
		t.Error("stable version should NOT be alpha")
	}
}

func TestIsAlpha_BetaNotAlpha(t *testing.T) {
	v := NewVersion("1.0.0-beta1")
	if v.IsAlpha() {
		t.Error("beta version should NOT be alpha")
	}
}

func TestIsAlpha_RCNotAlpha(t *testing.T) {
	v := NewVersion("1.0.0-rc1")
	if v.IsAlpha() {
		t.Error("rc version should NOT be alpha")
	}
}

// ─── IsBeta ───────────────────────────────────────────────────────────────────

func TestIsBeta_Beta(t *testing.T) {
	v := NewVersion("1.0.0-beta2")
	if !v.IsBeta() {
		t.Error("expected IsBeta() true")
	}
}

func TestIsBeta_StableNotBeta(t *testing.T) {
	v := NewVersion("1.0.0")
	if v.IsBeta() {
		t.Error("stable version should NOT be beta")
	}
}

func TestIsBeta_AlphaNotBeta(t *testing.T) {
	v := NewVersion("1.0.0-alpha1")
	if v.IsBeta() {
		t.Error("alpha version should NOT be beta")
	}
}

func TestIsBeta_RCNotBeta(t *testing.T) {
	v := NewVersion("1.0.0-rc1")
	if v.IsBeta() {
		t.Error("rc version should NOT be beta")
	}
}

// ─── IsRC ─────────────────────────────────────────────────────────────────────

func TestIsRC_RC(t *testing.T) {
	v := NewVersion("1.0.0-rc1")
	if !v.IsRC() {
		t.Error("expected IsRC() true")
	}
}

func TestIsRC_StableNotRC(t *testing.T) {
	v := NewVersion("1.0.0")
	if v.IsRC() {
		t.Error("stable version should NOT be rc")
	}
}

func TestIsRC_BetaNotRC(t *testing.T) {
	v := NewVersion("1.0.0-beta1")
	if v.IsRC() {
		t.Error("beta version should NOT be rc")
	}
}

func TestIsRC_AlphaNotRC(t *testing.T) {
	v := NewVersion("1.0.0-alpha1")
	if v.IsRC() {
		t.Error("alpha version should NOT be rc")
	}
}

// ─── IsPreview ────────────────────────────────────────────────────────────────

func TestIsPreview_Preview(t *testing.T) {
	v := NewVersion("1.0.0-preview.1")
	if !v.IsPreview() {
		t.Error("expected IsPreview() true")
	}
}

func TestIsPreview_StableNotPreview(t *testing.T) {
	v := NewVersion("1.0.0")
	if v.IsPreview() {
		t.Error("stable version should NOT be preview")
	}
}

func TestIsPreview_BetaNotPreview(t *testing.T) {
	v := NewVersion("1.0.0-beta1")
	if v.IsPreview() {
		t.Error("beta version should NOT be preview")
	}
}

func TestIsPreview_RCNotPreview(t *testing.T) {
	v := NewVersion("1.0.0-rc1")
	if v.IsPreview() {
		t.Error("rc version should NOT be preview")
	}
}

// ─── IsPseudo ─────────────────────────────────────────────────────────────────

func TestIsPseudo_WithBothFields(t *testing.T) {
	n := int64(5)
	v := Version{
		Timestamp:  []TimestampInfo{{Original: "20210101", Parsed: "2021-01-01T00:00:00Z"}},
		CommitHash: []CommitHashInfo{{Original: "abc123def456", Parsed: "abc123def456", CommitsSinceTag: &n}},
	}
	if !v.IsPseudo() {
		t.Error("expected IsPseudo() true when both Timestamp and CommitHash are set")
	}
}

func TestIsPseudo_OnlyTimestamp(t *testing.T) {
	v := Version{
		Timestamp:  []TimestampInfo{{Original: "20210101", Parsed: "2021-01-01T00:00:00Z"}},
		CommitHash: []CommitHashInfo{},
	}
	if v.IsPseudo() {
		t.Error("IsPseudo() should be false with only Timestamp")
	}
}

func TestIsPseudo_OnlyCommit(t *testing.T) {
	v := Version{
		Timestamp:  []TimestampInfo{},
		CommitHash: []CommitHashInfo{{Original: "abc123", Parsed: "abc123"}},
	}
	if v.IsPseudo() {
		t.Error("IsPseudo() should be false with only CommitHash")
	}
}

func TestIsPseudo_Neither(t *testing.T) {
	v := NewVersion("1.0.0")
	if v.IsPseudo() {
		t.Error("stable version should NOT be pseudo")
	}
}

// ─── SortVersions ─────────────────────────────────────────────────────────────

func TestSortVersions_Ascending(t *testing.T) {
	vs := parseVersions([]string{"3.0.0", "1.0.0", "2.0.0"})
	SortVersions(vs)
	assertOrder(t, vs, []string{"1.0.0", "2.0.0", "3.0.0"})
}

func TestSortVersions_Descending(t *testing.T) {
	vs := parseVersions([]string{"1.0.0", "3.0.0", "2.0.0"})
	SortVersions(vs, true)
	assertOrder(t, vs, []string{"3.0.0", "2.0.0", "1.0.0"})
}

func TestSortVersions_Empty(t *testing.T) {
	var vs []*Version
	SortVersions(vs)
	if len(vs) != 0 {
		t.Errorf("sorting empty slice should remain empty")
	}
}

func TestSortVersions_SingleElement(t *testing.T) {
	vs := parseVersions([]string{"1.0.0"})
	SortVersions(vs)
	if len(vs) != 1 || vs[0].String() != "1.0.0" {
		t.Error("single element sort should be unchanged")
	}
}

func TestSortVersions_PreReleaseBeforeStable(t *testing.T) {
	vs := parseVersions([]string{"1.0.0", "1.0.0-rc1", "1.0.0-alpha1", "1.0.0-beta1"})
	SortVersions(vs)
	// ascending: alpha < beta < rc < stable (last element should be stable)
	last := vs[len(vs)-1]
	if last.String() != "1.0.0" {
		t.Errorf("last element (stable) should be 1.0.0, got %q", last.String())
	}
	// first element should be alpha
	first := vs[0]
	if !first.IsAlpha() {
		t.Errorf("first element should be alpha, got %q", first.String())
	}
}

func TestSortVersions_MixedMajorVersions(t *testing.T) {
	vs := parseVersions([]string{"10.0.0", "9.0.0", "8.0.0", "11.0.0"})
	SortVersions(vs)
	assertOrder(t, vs, []string{"8.0.0", "9.0.0", "10.0.0", "11.0.0"})
}

func TestSortVersions_NoDescendingArg(t *testing.T) {
	vs := parseVersions([]string{"2.0.0", "1.0.0"})
	SortVersions(vs) // no descending arg
	if vs[0].String() != "1.0.0" {
		t.Error("default (no arg) should sort ascending")
	}
}

func TestSortVersions_FalseDescending(t *testing.T) {
	vs := parseVersions([]string{"2.0.0", "1.0.0"})
	SortVersions(vs, false) // explicitly ascending
	if vs[0].String() != "1.0.0" {
		t.Error("descending=false should sort ascending")
	}
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func assertCore(t *testing.T, v Version, major, minor, patch int64) {
	t.Helper()
	if safeInt(v.Major) != major {
		t.Errorf("major: got %d, want %d", safeInt(v.Major), major)
	}
	if safeInt(v.Minor) != minor {
		t.Errorf("minor: got %d, want %d", safeInt(v.Minor), minor)
	}
	if safeInt(v.Patch) != patch {
		t.Errorf("patch: got %d, want %d", safeInt(v.Patch), patch)
	}
}

func parseVersions(strs []string) []*Version {
	vs := make([]*Version, len(strs))
	for i, s := range strs {
		v := NewVersion(s)
		vs[i] = &v
	}
	return vs
}

func assertOrder(t *testing.T, vs []*Version, want []string) {
	t.Helper()
	if len(vs) != len(want) {
		t.Fatalf("length mismatch: got %d, want %d", len(vs), len(want))
	}
	for i, w := range want {
		got := vs[i].String()
		if got != w {
			t.Errorf("pos %d: got %q, want %q", i, got, w)
		}
	}
}

