package resolver

import (
	"testing"

	"github.com/rng70/versions/v2/vars"
)

func TestNPM_WildcardMinor(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNPM, "2.x", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{}) // no 2.x versions in TestVersions
}

func TestNPM_WildcardPatch(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNPM, "3.3.x", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{}) // no 3.3.x versions in TestVersions
}

func TestNPM_ExactVersion(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNPM, "2.0.1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{}) // no 2.0.1 in TestVersions
}

func TestNPM_TildeMinor(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNPM, "~1.2", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{}) // no 1.2.x versions
}

func TestNPM_TildePatch(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNPM, "~1.2.3", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{}) // no 1.2.x versions
}

func TestNPM_HyphenRange(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNPM, "1.0.0 - 2.9999.9999", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{}) // no versions in 1.0-2.9999 range
}

func TestNPM_GteAndLt(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNPM, ">=1.0.2 <2.1.2", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 84) // all TestVersions >= 1.0.2
}

func TestNPM_GtAndLte(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNPM, ">1.0.2 <=2.3.4", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 84)
}

func TestNPM_OrRanges(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNPM, "<1.0.0 || >=2.3.1 <2.4.5 || >=2.5.2 <3.0.0", vars.TestVersions)
	assertParsedCount(t, a, 3)
	assertMatchCount(t, a, 84)
}

func TestNPM_Latest(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNPM, "latest", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{}) // "latest" not in TestVersions
}

func TestNPM_PkgAlias(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNPM, "npm:pkg@1.0.0", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{}) // 1.0.0 not in TestVersions
}

func TestNPM_URLIgnored(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNPM, "http://npmjs.com/example.tar.gz", vars.TestVersions)
	assertParsedCount(t, a, 0)
	assertMatches(t, a, []string{})
}

func TestNPM_AlphaFinalRange(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNPM, ">=4.2.0.Alpha1, <4.2.3.Final", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{})
}

func TestNPM_PreRelease(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNPM, ">= 9.0.0-preview.1.24081.5, <= 9.0.0-rc.1.24452.1", preReleaseVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, wantPreRelease)
}

func TestNPM_Exact(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNPM, "= 10.0.1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1"})
}

func TestNPM_BareExact(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNPM, "10.0.1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1"})
}

func TestNPM_BareExactPreRelease(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNPM, "10.0.1-beta1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1-beta1"})
}

func TestNPM_Range(t *testing.T) {
	npmVersions := []string{
		"0.0.0-PLACEHOLDER", "17.0.0-rc.0", "17.0.0", "17.0.1", "17.0.10",
		"18.0.0-rc.0", "18.0.0", "18.2.21",
		"19.0.0-rc.0", "19.0.0-next.0", "19.0.0-next.13", "19.0.0", "19.0.7",
		"19.1.0-rc.0", "19.1.0", "19.1.9", "19.2.0-rc.0", "19.2.0", "19.2.17",
		"19.2.18", "20.0.0-rc.0", "20.0.0", "20.3.6",
		"21.0.0-next.0", "21.0.0-next.7",
	}

	a := AnalyzeConstraint(vars.StyleNPM, ">=19.0.0-next.0, <19.2.18", npmVersions)
	assertParsedCount(t, a, 1)
	if len(a.Matches) == 0 {
		t.Error("expected matches for >=19.0.0-next.0, <19.2.18")
	}
	for _, m := range a.Matches {
		if m == "19.2.18" {
			t.Error("19.2.18 should NOT match <19.2.18")
		}
	}
}
