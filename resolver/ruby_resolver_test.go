package resolver

import (
	"testing"

	"github.com/rng70/versions/vars"
)

func TestRuby_PessimisticMinor(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRuby, "~> 10.0", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1", "10.0.1-beta1", "10.0.2", "10.0.3"})
}

func TestRuby_PessimisticPatch(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRuby, "~> 10.0.1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1", "10.0.2", "10.0.3"})
}

func TestRuby_Gte(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRuby, ">= 10.0.1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 30) // 10.0.1+ (excludes beta of 10.0.1)
}

func TestRuby_Lt(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRuby, "< 12.0.0", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 62)
}

func TestRuby_Exact(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRuby, "= 10.0.2", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.2"})
}

func TestRuby_NotEqual(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRuby, "!= 10.0.1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 83) // all except 10.0.1
}

func TestRuby_Compound(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRuby, ">= 10.0.0, < 12.0.0", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 9)
}

func TestRuby_PreRelease(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRuby, ">= 9.0.0-preview.1.24081.5, <= 9.0.0-rc.1.24452.1", preReleaseVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, wantPreRelease)
}

func TestRuby_BareExact(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRuby, "10.0.1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1"})
}

func TestRuby_BareExactPreRelease(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRuby, "10.0.1-beta1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1-beta1"})
}
