package resolver

import (
	"testing"

	"github.com/rng70/versions/v2/vars"
)

func TestGo_GteWithPrefix(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleGo, ">= v10.0.1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 30)
}

func TestGo_LtWithPrefix(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleGo, "< v12.0.0", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 62)
}

func TestGo_RangeWithPrefix(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleGo, ">= v10.0.0, < v12.0.0", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 9)
}

func TestGo_ExactWithPrefix(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleGo, "v10.0.1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1"})
}

func TestGo_GteNoPrefix(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleGo, ">= 10.0.0", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 31)
}

func TestGo_NotEqual(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleGo, "!= v10.0.1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 83)
}

func TestGo_PreRelease(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleGo, ">= v9.0.0-preview.1.24081.5, <= v9.0.0-rc.1.24452.1", preReleaseVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, wantPreRelease)
}

func TestGo_Exact(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleGo, "= v10.0.1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1"})
}

func TestGo_BareExact(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleGo, "10.0.1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1"})
}

func TestGo_BareExactPreRelease(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleGo, "10.0.1-beta1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1-beta1"})
}
