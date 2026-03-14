package resolver

import (
	"testing"

	"github.com/rng70/versions/vars"
)

func TestPython_CompatibleRelease(t *testing.T) {
	a := AnalyzeConstraint(vars.StylePy, "~=1.4", PyPITestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{})
}

func TestPython_WildcardEqual(t *testing.T) {
	a := AnalyzeConstraint(vars.StylePy, "==1.2.*", PyPITestVersions)
	assertParsedCount(t, a, 0) // fails to parse (wildcard with ==)
	assertMatches(t, a, []string{})
}

func TestPython_ExactVersion(t *testing.T) {
	a := AnalyzeConstraint(vars.StylePy, "==2.0.1", PyPITestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{})
}

func TestPython_CompatibleReleasePatch(t *testing.T) {
	a := AnalyzeConstraint(vars.StylePy, "~=1.2.3", PyPITestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{})
}

func TestPython_GteAndLt(t *testing.T) {
	a := AnalyzeConstraint(vars.StylePy, ">=1.0.2, <2.1.2", PyPITestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"1.0.2", "1.0.3", "2.0.0beta", "2.0.3"})
}

func TestPython_GtAndLte(t *testing.T) {
	a := AnalyzeConstraint(vars.StylePy, ">1.0.2, <=2.3.4", PyPITestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"1.0.3", "2.0.0beta", "2.0.3", "2.1.2", "2.1.3", "2.3.0", "2.3.1", "2.3.4"})
}

func TestPython_NotEqual(t *testing.T) {
	a := AnalyzeConstraint(vars.StylePy, "!=2.3.1, >=1.0.0, <3.0.0", PyPITestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{
		"1.0.0",
		"1.0.2",
		"1.0.3",
		"2.0.0beta",
		"2.0.3",
		"2.1.2",
		"2.1.3",
		"2.3.0",
		"2.3.4",
		"2.3.5",
		"2.4.1",
		"2.5.0",
		"2.6.0",
		"2.6.1",
		"3.0.0a2",
		"3.0.0a3",
		"3.0.0b1",
		"3.0.0b1.post1",
		"3.0.0b1.post2",
		"3.0.0b2",
		"3.0.0b2.post1",
		"3.0.0b2.post2",
		"3.0.0b3",
		"3.0.0b4",
	})
}

func TestPython_BareVersion(t *testing.T) {
	a := AnalyzeConstraint(vars.StylePy, "2.3.1", PyPITestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"2.3.1"})
}

func TestPython_BareVersion_Missing(t *testing.T) {
	a := AnalyzeConstraint(vars.StylePy, "2.3.2", PyPITestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{})
}

func TestPython_FourPartVersion(t *testing.T) {
	a := AnalyzeConstraint(vars.StylePy, ">=1.2.0.1", PyPITestVersions)
	assertParsedCount(t, a, 0)
	assertMatches(t, a, []string{})
}

func TestPython_RCVersion(t *testing.T) {
	a := AnalyzeConstraint(vars.StylePy, ">= 6.30.0rc1, <= 6.33.4", PyPITestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{
		"6.30.0",
		"6.30.0rc1",
		"6.30.0rc2",
		"6.30.1",
		"6.30.2",
		"6.31.0",
		"6.31.0rc1",
		"6.31.0rc2",
		"6.31.1",
		"6.32.0",
		"6.32.0rc1",
		"6.32.0rc2",
		"6.32.1",
		"6.33.0",
		"6.33.0rc1",
		"6.33.0rc2",
		"6.33.1",
		"6.33.2",
		"6.33.3",
		"6.33.4",
	})
}

func TestPython_LtOnly(t *testing.T) {
	a := AnalyzeConstraint(vars.StylePy, "< 5.29.6", PyPITestVersions)
	assertParsedCount(t, a, 1)
	assertMatchCount(t, a, 192)
	assertMatches(t, a, []string{
		"1.0.0",
		"1.0.2",
		"1.0.3",
		"2.0.0beta",
		"2.0.3",
		"2.1.2",
		"2.1.3",
		"2.3.0",
		"2.3.1",
		"2.3.4",
		"2.3.5",
		"2.4.1",
		"2.5.0",
		"2.6.0",
		"2.6.1",
		"3.0.0",
		"3.0.0a2",
		"3.0.0a3",
		"3.0.0b1",
		"3.0.0b1.post1",
		"3.0.0b1.post2",
		"3.0.0b2",
		"3.0.0b2.post1",
		"3.0.0b2.post2",
		"3.0.0b3",
		"3.0.0b4",
		"3.1.0",
		"3.10.0",
		"3.10.0rc1",
		"3.1.0.post1",
		"3.11.0",
		"3.11.0rc1",
		"3.11.0rc2",
		"3.11.1",
		"3.11.2",
		"3.11.3",
		"3.12.0",
		"3.12.0rc1",
		"3.12.0rc2",
		"3.12.1",
		"3.12.2",
		"3.12.4",
		"3.13.0",
		"3.13.0rc3",
		"3.14.0",
		"3.14.0rc1",
		"3.14.0rc2",
		"3.14.0rc3",
		"3.15.0",
		"3.15.0rc1",
		"3.15.0rc2",
		"3.15.1",
		"3.15.2",
		"3.15.3",
		"3.15.4",
		"3.15.5",
		"3.15.6",
		"3.15.7",
		"3.15.8",
		"3.16.0",
		"3.16.0rc1",
		"3.16.0rc2",
		"3.17.0",
		"3.17.0rc1",
		"3.17.0rc2",
		"3.17.1",
		"3.17.2",
		"3.17.3",
		"3.18.0",
		"3.18.0rc1",
		"3.18.0rc2",
		"3.18.1",
		"3.18.3",
		"3.19.0",
		"3.19.0rc1",
		"3.19.0rc2",
		"3.19.1",
		"3.19.2",
		"3.19.3",
		"3.19.4",
		"3.19.5",
		"3.19.6",
		"3.2.0",
		"3.20.0",
		"3.20.0rc1",
		"3.20.0rc2",
		"3.20.1",
		"3.20.1rc1",
		"3.20.2",
		"3.20.3",
		"3.2.0rc1",
		"3.2.0rc1.post1",
		"3.2.0rc2",
		"3.3.0",
		"3.4.0",
		"3.5.0.post1",
		"3.5.1",
		"3.5.2",
		"3.5.2.post1",
		"3.6.0",
		"3.6.1",
		"3.7.0",
		"3.7.0rc2",
		"3.7.0rc3",
		"3.7.1",
		"3.8.0",
		"3.8.0rc1",
		"3.9.0",
		"3.9.0rc1",
		"3.9.1",
		"3.9.2",
		"4.0.0rc1",
		"4.0.0rc2",
		"4.21.0",
		"4.21.0rc1",
		"4.21.0rc2",
		"4.21.1",
		"4.21.10",
		"4.21.11",
		"4.21.12",
		"4.21.2",
		"4.21.3",
		"4.21.4",
		"4.21.5",
		"4.21.6",
		"4.21.7",
		"4.21.8",
		"4.21.9",
		"4.22.0",
		"4.22.0rc1",
		"4.22.0rc2",
		"4.22.0rc3",
		"4.22.1",
		"4.22.3",
		"4.22.4",
		"4.22.5",
		"4.23.0",
		"4.23.0rc2",
		"4.23.0rc3",
		"4.23.1",
		"4.23.2",
		"4.23.3",
		"4.23.4",
		"4.24.0",
		"4.24.0rc1",
		"4.24.0rc2",
		"4.24.0rc3",
		"4.24.1",
		"4.24.2",
		"4.24.3",
		"4.24.4",
		"4.25.0",
		"4.25.0rc1",
		"4.25.0rc2",
		"4.25.1",
		"4.25.2",
		"4.25.3",
		"4.25.4",
		"4.25.5",
		"4.25.6",
		"4.25.7",
		"4.25.8",
		"5.26.0",
		"5.26.0rc1",
		"5.26.0rc2",
		"5.26.0rc3",
		"5.26.1",
		"5.27.0",
		"5.27.0rc1",
		"5.27.0rc2",
		"5.27.0rc3",
		"5.27.1",
		"5.27.2",
		"5.27.3",
		"5.27.4",
		"5.27.5",
		"5.28.0",
		"5.28.0rc1",
		"5.28.0rc2",
		"5.28.0rc3",
		"5.28.1",
		"5.28.2",
		"5.28.3",
		"5.29.0",
		"5.29.0rc1",
		"5.29.0rc2",
		"5.29.0rc3",
		"5.29.1",
		"5.29.2",
		"5.29.3",
		"5.29.4",
		"5.29.5",
	})
}

func TestPython_PreRelease(t *testing.T) {
	a := AnalyzeConstraint(vars.StylePy, ">= 9.0.0-preview.1.24081.5, <= 9.0.0-rc.1.24452.1", preReleaseVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, wantPreRelease)
}

func TestPython_BareExact(t *testing.T) {
	a := AnalyzeConstraint(vars.StylePy, "3.0.0b1.post1", PyPITestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"3.0.0b1.post1"})
}

func TestPython_BareExact_Missing(t *testing.T) {
	a := AnalyzeConstraint(vars.StylePy, "10.0.1", PyPITestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{})
}

func TestPython_BareExactPreRelease(t *testing.T) {
	a := AnalyzeConstraint(vars.StylePy, "7.35.0-beta1", PyPITestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{"7.35.0-beta1"})
}

func TestPython_BareExactPreRelease_Missing(t *testing.T) {
	a := AnalyzeConstraint(vars.StylePy, "10.0.1-beta1", PyPITestVersions)
	assertParsedCount(t, a, 1)
	assertMatches(t, a, []string{})
}

func TestPython_Range(t *testing.T) {
	pyVersions := []string{
		"1.1.1", "3.0.0", "2.9.9", "2.9.0", "1.9.0", "2.8.1",
		"1.0.0", "2.0.0", "1.2.0", "1.0.0-alpha", "1.0.0-beta.1",
	}

	a := AnalyzeConstraint(vars.StylePy, ">1.1.0,<2.9.9", pyVersions)
	assertParsedCount(t, a, 1)
	if len(a.Matches) == 0 {
		t.Error("expected matches for >1.1.0,<2.9.9")
	}
	for _, m := range a.Matches {
		if m == "2.9.9" || m == "3.0.0" {
			t.Errorf("%s should NOT match >1.1.0,<2.9.9", m)
		}
	}
}

var (
	PyPITestVersions = []string{
		"1.0.0",
		"1.0.2",
		"1.0.3",
		"2.0.0beta",
		"2.0.3",
		"2.1.2",
		"2.1.3",
		"2.3.0",
		"2.3.1",
		"2.3.4",
		"2.3.5",
		"2.4.1",
		"2.5.0",
		"2.6.0",
		"2.6.1",
		"3.0.0",
		"3.0.0a2",
		"3.0.0a3",
		"3.0.0b1",
		"3.0.0b1.post1",
		"3.0.0b1.post2",
		"3.0.0b2",
		"3.0.0b2.post1",
		"3.0.0b2.post2",
		"3.0.0b3",
		"3.0.0b4",
		"3.1.0",
		"3.10.0",
		"3.10.0rc1",
		"3.1.0.post1",
		"3.11.0",
		"3.11.0rc1",
		"3.11.0rc2",
		"3.11.1",
		"3.11.2",
		"3.11.3",
		"3.12.0",
		"3.12.0rc1",
		"3.12.0rc2",
		"3.12.1",
		"3.12.2",
		"3.12.4",
		"3.13.0",
		"3.13.0rc3",
		"3.14.0",
		"3.14.0rc1",
		"3.14.0rc2",
		"3.14.0rc3",
		"3.15.0",
		"3.15.0rc1",
		"3.15.0rc2",
		"3.15.1",
		"3.15.2",
		"3.15.3",
		"3.15.4",
		"3.15.5",
		"3.15.6",
		"3.15.7",
		"3.15.8",
		"3.16.0",
		"3.16.0rc1",
		"3.16.0rc2",
		"3.17.0",
		"3.17.0rc1",
		"3.17.0rc2",
		"3.17.1",
		"3.17.2",
		"3.17.3",
		"3.18.0",
		"3.18.0rc1",
		"3.18.0rc2",
		"3.18.1",
		"3.18.3",
		"3.19.0",
		"3.19.0rc1",
		"3.19.0rc2",
		"3.19.1",
		"3.19.2",
		"3.19.3",
		"3.19.4",
		"3.19.5",
		"3.19.6",
		"3.2.0",
		"3.20.0",
		"3.20.0rc1",
		"3.20.0rc2",
		"3.20.1",
		"3.20.1rc1",
		"3.20.2",
		"3.20.3",
		"3.2.0rc1",
		"3.2.0rc1.post1",
		"3.2.0rc2",
		"3.3.0",
		"3.4.0",
		"3.5.0.post1",
		"3.5.1",
		"3.5.2",
		"3.5.2.post1",
		"3.6.0",
		"3.6.1",
		"3.7.0",
		"3.7.0rc2",
		"3.7.0rc3",
		"3.7.1",
		"3.8.0",
		"3.8.0rc1",
		"3.9.0",
		"3.9.0rc1",
		"3.9.1",
		"3.9.2",
		"4.0.0rc1",
		"4.0.0rc2",
		"4.21.0",
		"4.21.0rc1",
		"4.21.0rc2",
		"4.21.1",
		"4.21.10",
		"4.21.11",
		"4.21.12",
		"4.21.2",
		"4.21.3",
		"4.21.4",
		"4.21.5",
		"4.21.6",
		"4.21.7",
		"4.21.8",
		"4.21.9",
		"4.22.0",
		"4.22.0rc1",
		"4.22.0rc2",
		"4.22.0rc3",
		"4.22.1",
		"4.22.3",
		"4.22.4",
		"4.22.5",
		"4.23.0",
		"4.23.0rc2",
		"4.23.0rc3",
		"4.23.1",
		"4.23.2",
		"4.23.3",
		"4.23.4",
		"4.24.0",
		"4.24.0rc1",
		"4.24.0rc2",
		"4.24.0rc3",
		"4.24.1",
		"4.24.2",
		"4.24.3",
		"4.24.4",
		"4.25.0",
		"4.25.0rc1",
		"4.25.0rc2",
		"4.25.1",
		"4.25.2",
		"4.25.3",
		"4.25.4",
		"4.25.5",
		"4.25.6",
		"4.25.7",
		"4.25.8",
		"5.26.0",
		"5.26.0rc1",
		"5.26.0rc2",
		"5.26.0rc3",
		"5.26.1",
		"5.27.0",
		"5.27.0rc1",
		"5.27.0rc2",
		"5.27.0rc3",
		"5.27.1",
		"5.27.2",
		"5.27.3",
		"5.27.4",
		"5.27.5",
		"5.28.0",
		"5.28.0rc1",
		"5.28.0rc2",
		"5.28.0rc3",
		"5.28.1",
		"5.28.2",
		"5.28.3",
		"5.29.0",
		"5.29.0rc1",
		"5.29.0rc2",
		"5.29.0rc3",
		"5.29.1",
		"5.29.2",
		"5.29.3",
		"5.29.4",
		"5.29.5",
		"6.30.0",
		"6.30.0rc1",
		"6.30.0rc2",
		"6.30.1",
		"6.30.2",
		"6.31.0",
		"6.31.0rc1",
		"6.31.0rc2",
		"6.31.1",
		"6.32.0",
		"6.32.0rc1",
		"6.32.0rc2",
		"6.32.1",
		"6.33.0",
		"6.33.0rc1",
		"6.33.0rc2",
		"6.33.1",
		"6.33.2",
		"6.33.3",
		"6.33.4",
		"6.33.5",
		"7.34.0rc1",
		"7.35.0-beta1",
	}
)
