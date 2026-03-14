package resolver

import (
	"testing"

	"github.com/rng70/versions/vars"
)

func TestNuGet_ClosedOpen_Null(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "[1.0.0, 2.0.0)", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{})
}

func TestNuGet_ClosedOpen(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "[9.0.0-preview.1.24081.5, 10.0.0)", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{
		"9.0.0",
		"9.0.0-preview.1.24081.5",
		"9.0.0-preview.2.24128.4",
		"9.0.0-preview.3.24172.13",
		"9.0.0-preview.4.24267.6",
		"9.0.0-preview.5.24306.11",
		"9.0.0-preview.6.24328.4",
		"9.0.0-preview.7.24406.2",
		"9.0.0-rc.1.24452.1",
		"9.0.0-rc.2.24474.3",
		"9.0.1",
		"9.0.10",
		"9.0.11",
		"9.0.12",
		"9.0.13",
		"9.0.14",
		"9.0.2",
		"9.0.3",
		"9.0.4",
		"9.0.5",
		"9.0.6",
		"9.0.7",
		"9.0.8",
		"9.0.9",
		"10.0.0-preview.1.25120.3",
		"10.0.0-preview.2.25164.1",
		"10.0.0-preview.3.25172.1",
		"10.0.0-preview.4.25258.110",
		"10.0.0-preview.5.25277.114",
		"10.0.0-preview.6.25358.103",
		"10.0.0-preview.7.25380.108",
		"10.0.0-rc.1.25451.107",
		"10.0.0-rc.2.25502.107",
	})
}

func TestNuGet_OpenClosed_Null(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "(,2.4.5]", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{})
}

func TestNuGet_OpenClosed(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "(,8.0.0]", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{
		// 6.x
		"6.0.0",
		"6.0.0-preview.1.21103.6",
		"6.0.0-preview.2.21154.6",
		"6.0.0-preview.3.21201.13",
		"6.0.0-preview.4.21253.5",
		"6.0.0-preview.5.21301.17",
		"6.0.0-preview.6.21355.2",
		"6.0.0-preview.7.21378.6",
		"6.0.0-rc.1.21452.15",
		"6.0.0-rc.2.21480.10",
		"6.0.1",
		"6.0.10",
		"6.0.11",
		"6.0.12",
		"6.0.13",
		"6.0.14",
		"6.0.15",
		"6.0.16",
		"6.0.18",
		"6.0.19",
		"6.0.2",
		"6.0.20",
		"6.0.21",
		"6.0.22",
		"6.0.23",
		"6.0.24",
		"6.0.25",
		"6.0.26",
		"6.0.27",
		"6.0.28",
		"6.0.29",
		"6.0.3",
		"6.0.30",
		"6.0.31",
		"6.0.32",
		"6.0.33",
		"6.0.35",
		"6.0.36",
		"6.0.4",
		"6.0.5",
		"6.0.6",
		"6.0.7",
		"6.0.8",
		"6.0.9",
		// 7.x
		"7.0.0",
		"7.0.0-preview.1.22109.13",
		"7.0.0-preview.2.22153.2",
		"7.0.0-preview.3.22178.4",
		"7.0.0-preview.4.22251.1",
		"7.0.0-preview.5.22303.8",
		"7.0.0-preview.6.22330.3",
		"7.0.0-preview.7.22376.6",
		"7.0.0-rc.1.22427.2",
		"7.0.0-rc.2.22476.2",
		"7.0.1",
		"7.0.10",
		"7.0.11",
		"7.0.12",
		"7.0.13",
		"7.0.14",
		"7.0.15",
		"7.0.16",
		"7.0.17",
		"7.0.18",
		"7.0.19",
		"7.0.2",
		"7.0.20",
		"7.0.3",
		"7.0.4",
		"7.0.5",
		"7.0.7",
		"7.0.8",
		"7.0.9",
		// 8.x
		"8.0.0",
		"8.0.0-preview.1.23112.2",
		"8.0.0-preview.2.23153.2",
		"8.0.0-preview.3.23177.8",
		"8.0.0-preview.4.23260.4",
		"8.0.0-preview.5.23302.2",
		"8.0.0-preview.6.23329.11",
		"8.0.0-preview.7.23375.9",
		"8.0.0-rc.1.23421.29",
		"8.0.0-rc.2.23480.2",
	})
}

func TestNuGet_ClosedUnbounded(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "[2.5.2,)", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 162) // all versions >= 2.5.2
}

func TestNuGet_ExactBracket(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "[2.0.1]", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{}) // no 2.0.1 in NugetTestVersions
}

func TestNuGet_WildcardMinor(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "1.*", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{})
}

func TestNuGet_WildcardPatch(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "1.2.*", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{})
}

func TestNuGet_ExactVersion(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "2.3.1", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{})
}

func TestNuGet_LtOnly(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "< 13.0.1", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 162)
}

func TestNuGet_GtAndLt(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, ">11.0.1, <13.0.1", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 13)
}

func TestNuGet_PreRelease(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, ">= 9.0.0-preview.1.24081.5, <= 9.0.0-rc.1.24452.1", preReleaseVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, wantPreRelease)
}

func TestNuGet_BareExact(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "10.0.1", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1"})
}

func TestNuGet_BareExactPreRelease(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "10.0.1-beta1", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"10.0.1-beta1"})
}

// ─── New test cases ───────────────────────────────────────────────────────────

// TestNuGet_GteAndLt uses operator syntax for >= 8.0.0, < 9.0.0.
// Matches: 8.0.0 stable (1) + 8.0.1–8.0.25 (24) + 9.0.0-preview.x (7) + 9.0.0-rc.x (2) = 34.
// Note: 9.x pre-releases satisfy < 9.0.0 because pre-release is less than the stable release.
func TestNuGet_GteAndLt(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, ">= 8.0.0, < 9.0.0", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 34)
}

// TestNuGet_ClosedOpen_Present uses bracket notation [8.0.0, 9.0.0) — equivalent to
// >= 8.0.0 and < 9.0.0. Should produce the same 34 matches as TestNuGet_GteAndLt.
func TestNuGet_ClosedOpen_Present(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "[8.0.0, 9.0.0)", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 34)
}

// TestNuGet_OpenOpen uses (8.0.0, 9.0.0) — strictly between the two bounds.
// Excludes 8.0.0 stable but still captures 9.x pre-releases.
// Matches: 8.0.1–8.0.25 (24) + 9.0.0-preview.x (7) + 9.0.0-rc.x (2) = 33.
func TestNuGet_OpenOpen(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "(8.0.0, 9.0.0)", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 33)
}

// TestNuGet_WildcardMajor_Present tests 6.* which expands to >= 6.0.0, < 7.0.0.
// Matches: 6.0.0 stable (1) + 6.0.1–6.0.36 (34) + 7.0.0-preview.x (7) + 7.0.0-rc.x (2) = 44.
// 7.x pre-releases have core 7.0.0 > 6.0.0 and are < 7.0.0 (stable), so they fall in range.
func TestNuGet_WildcardMajor_Present(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "6.*", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 44)
}

// TestNuGet_WildcardPatch_Present tests 7.0.* which expands to >= 7.0.0, < 7.1.0.
// Matches: 7.0.0 stable (1) + 7.0.1–7.0.20 (19) = 20.
// 7.0.x pre-releases are excluded (< 7.0.0); 8.x pre-releases have core 8.0.0 > 7.1.0 so excluded.
func TestNuGet_WildcardPatch_Present(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "7.0.*", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 20)
}

// TestNuGet_LteOnly tests <= 7.0.0.
// Matches: all 6.x (44) + 7.0.0-preview.x (7) + 7.0.0-rc.x (2) + 7.0.0 stable (1) = 54.
func TestNuGet_LteOnly(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "<= 7.0.0", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 54)
}

// TestNuGet_GteOnly tests >= 9.0.0.
// 9.x pre-releases are excluded (< 9.0.0 stable).
// Matches: 9.0.0 (1) + 9.0.1–9.0.14 (14) + all 10.x (15) + all 11.x (9) + all 12.x (7) = 46.
func TestNuGet_GteOnly(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, ">= 9.0.0", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 46)
}

// TestNuGet_GtOnly tests > 10.0.0 (strictly greater).
// Excludes 10.0.0 stable and all 10.0.0 pre-releases.
// Matches: 10.0.1 (1) + 10.0.1-beta1 (1) + 10.0.2–10.0.4 (3) + all 11.x (9) + all 12.x (7) = 21.
func TestNuGet_GtOnly(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "> 10.0.0", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 21)
}

// TestNuGet_ExactBracket_Present tests [6.0.5] which matches exactly one version.
func TestNuGet_ExactBracket_Present(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "[6.0.5]", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"6.0.5"})
}

// TestNuGet_ExactVersion_Present tests a bare version that is present in the list.
func TestNuGet_ExactVersion_Present(t *testing.T) {
	a := AnalyzeConstraint(vars.StyleNuGet, "8.0.1", NugetTestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"8.0.1"})
}

var (
	NugetTestVersions = []string{
		// 10.x
		"10.0.0",
		"10.0.0-preview.1.25120.3",
		"10.0.0-preview.2.25164.1",
		"10.0.0-preview.3.25172.1",
		"10.0.0-preview.4.25258.110",
		"10.0.0-preview.5.25277.114",
		"10.0.0-preview.6.25358.103",
		"10.0.0-preview.7.25380.108",
		"10.0.0-rc.1.25451.107",
		"10.0.0-rc.2.25502.107",
		"10.0.1",
		"10.0.1-beta1",
		"10.0.2",
		"10.0.3",
		"10.0.4",
		// 11.x
		"11.0.0-preview.1.26104.118",
		"11.0.0-preview.2.26159.112",
		"11.0.1",
		"11.0.2",
		"11.0.3",
		"11.0.4",
		"11.0.5",
		"11.0.6",
		"11.0.7",
		// 12.x
		"12.0.0",
		"12.0.1",
		"12.0.2",
		"12.0.3",
		"12.0.4",
		"12.0.5",
		"12.0.6",
		// 6.x
		"6.0.0",
		"6.0.0-preview.1.21103.6",
		"6.0.0-preview.2.21154.6",
		"6.0.0-preview.3.21201.13",
		"6.0.0-preview.4.21253.5",
		"6.0.0-preview.5.21301.17",
		"6.0.0-preview.6.21355.2",
		"6.0.0-preview.7.21378.6",
		"6.0.0-rc.1.21452.15",
		"6.0.0-rc.2.21480.10",
		"6.0.1",
		"6.0.10",
		"6.0.11",
		"6.0.12",
		"6.0.13",
		"6.0.14",
		"6.0.15",
		"6.0.16",
		"6.0.18",
		"6.0.19",
		"6.0.2",
		"6.0.20",
		"6.0.21",
		"6.0.22",
		"6.0.23",
		"6.0.24",
		"6.0.25",
		"6.0.26",
		"6.0.27",
		"6.0.28",
		"6.0.29",
		"6.0.3",
		"6.0.30",
		"6.0.31",
		"6.0.32",
		"6.0.33",
		"6.0.35",
		"6.0.36",
		"6.0.4",
		"6.0.5",
		"6.0.6",
		"6.0.7",
		"6.0.8",
		"6.0.9",
		// 7.x
		"7.0.0",
		"7.0.0-preview.1.22109.13",
		"7.0.0-preview.2.22153.2",
		"7.0.0-preview.3.22178.4",
		"7.0.0-preview.4.22251.1",
		"7.0.0-preview.5.22303.8",
		"7.0.0-preview.6.22330.3",
		"7.0.0-preview.7.22376.6",
		"7.0.0-rc.1.22427.2",
		"7.0.0-rc.2.22476.2",
		"7.0.1",
		"7.0.10",
		"7.0.11",
		"7.0.12",
		"7.0.13",
		"7.0.14",
		"7.0.15",
		"7.0.16",
		"7.0.17",
		"7.0.18",
		"7.0.19",
		"7.0.2",
		"7.0.20",
		"7.0.3",
		"7.0.4",
		"7.0.5",
		"7.0.7",
		"7.0.8",
		"7.0.9",
		// 8.x
		"8.0.0",
		"8.0.0-preview.1.23112.2",
		"8.0.0-preview.2.23153.2",
		"8.0.0-preview.3.23177.8",
		"8.0.0-preview.4.23260.4",
		"8.0.0-preview.5.23302.2",
		"8.0.0-preview.6.23329.11",
		"8.0.0-preview.7.23375.9",
		"8.0.0-rc.1.23421.29",
		"8.0.0-rc.2.23480.2",
		"8.0.1",
		"8.0.10",
		"8.0.11",
		"8.0.12",
		"8.0.13",
		"8.0.14",
		"8.0.15",
		"8.0.16",
		"8.0.17",
		"8.0.18",
		"8.0.19",
		"8.0.2",
		"8.0.20",
		"8.0.21",
		"8.0.22",
		"8.0.23",
		"8.0.24",
		"8.0.25",
		"8.0.3",
		"8.0.4",
		"8.0.5",
		"8.0.6",
		"8.0.7",
		"8.0.8",
		// 9.x
		"9.0.0",
		"9.0.0-preview.1.24081.5",
		"9.0.0-preview.2.24128.4",
		"9.0.0-preview.3.24172.13",
		"9.0.0-preview.4.24267.6",
		"9.0.0-preview.5.24306.11",
		"9.0.0-preview.6.24328.4",
		"9.0.0-preview.7.24406.2",
		"9.0.0-rc.1.24452.1",
		"9.0.0-rc.2.24474.3",
		"9.0.1",
		"9.0.10",
		"9.0.11",
		"9.0.12",
		"9.0.13",
		"9.0.14",
		"9.0.2",
		"9.0.3",
		"9.0.4",
		"9.0.5",
		"9.0.6",
		"9.0.7",
		"9.0.8",
		"9.0.9",
	}
)
