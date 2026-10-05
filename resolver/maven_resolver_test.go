package resolver

import (
	"testing"

	"github.com/rng70/versions/v2/vars"
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

func TestMaven_Exact(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleMaven, "= 10.0.1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1"})
}

func TestMaven_QualifierUpperBound(t *testing.T) {
	// CVE-2024-47535 (io.netty:netty-common) constraint
	a := AnalyzeConstraint(vars.StyleMaven, "<= 4.1.114.Final", qualifiedVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"4.1.100.Final", "4.1.114.Final", "2.0.M1", "2.0.0"})
}

func TestMaven_QualifierRange(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleMaven, ">= 4.1.0.Final, <= 4.1.114.Final", qualifiedVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"4.1.100.Final", "4.1.114.Final"})
}

func TestMaven_QualifierExclusiveUpper(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleMaven, ">= 5.0.0, < 5.3.1.RELEASE", qualifiedVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"5.3.0.RELEASE"})
}

func TestMaven_QualifierExactBracket(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleMaven, "[4.1.114.Final]", qualifiedVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"4.1.114.Final"})
}

func TestMaven_QualifierBareExact(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleMaven, "4.1.114.Final", qualifiedVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"4.1.114.Final"})
}

func TestMaven_QualifierBeforePatchExact(t *testing.T) {
	// "2.0.M1" must not collapse to "2.0.0"
	a := AnalyzeConstraint(vars.StyleMaven, "= 2.0.M1", qualifiedVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"2.0.M1"})
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
