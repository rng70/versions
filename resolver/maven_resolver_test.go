package resolver

import (
	"testing"

	"github.com/rng70/versions/vars"
)

func TestMaven_BracketRange(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleMaven, "[10.0.0,12.0.0)", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 9) // 10.0.1..11.0.2 range
}

func TestMaven_ExactBracket(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleMaven, "[10.0.1]", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1"})
}

func TestMaven_UpperBoundInclusive(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleMaven, "(,11.0.1]", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 61)
}

func TestMaven_LowerBoundInclusive(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleMaven, "[11.0.0,)", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 27)
}

func TestMaven_SoftRequirement(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleMaven, "10.0.1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1"})
}

func TestMaven_ExclusiveRange(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleMaven, "(10.0.3,13.0.1)", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 17)
}

func TestMaven_PreRelease(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleMaven, ">= 9.0.0-preview.1.24081.5, <= 9.0.0-rc.1.24452.1", preReleaseVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, wantPreRelease)
}

func TestMaven_BareExact(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleMaven, "10.0.1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1"})
}

func TestMaven_BareExactPreRelease(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleMaven, "10.0.1-beta1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1-beta1"})
}
