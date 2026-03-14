package resolver

import (
	"testing"

	"github.com/rng70/versions/v2/vars"
)

func TestRust_Caret(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRust, "^10.0.1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1", "10.0.2", "10.0.3"})
}

func TestRust_Tilde(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRust, "~10.0.1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1", "10.0.2", "10.0.3"})
}

func TestRust_BareVersion(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRust, "10.0.1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1"})
}

func TestRust_Wildcard(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRust, "*", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, len(vars.TestVersions)) // all versions
}

func TestRust_WildcardMinor(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRust, "10.*", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1", "10.0.1-beta1", "10.0.2", "10.0.3"})
}

func TestRust_ExplicitRange(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRust, ">= 10.0.0, < 12.0.0", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 9)
}

func TestRust_Exact(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRust, "= 10.0.2", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.2"})
}

func TestRust_PreRelease(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRust, ">= 9.0.0-preview.1.24081.5, <= 9.0.0-rc.1.24452.1", preReleaseVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, wantPreRelease)
}

func TestRust_BareExact(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRust, "10.0.1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1"})
}

func TestRust_BareExactPreRelease(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleRust, "10.0.1-beta1", vars.TestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1-beta1"})
}
