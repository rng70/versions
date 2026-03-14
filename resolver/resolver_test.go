package resolver

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/rng70/versions/vars"
)

// assertParsedCount checks that the constraint string was parsed into the
// expected number of OR-groups.
func assertParsedCount(t *testing.T, a vars.Analysis, want int) {
	t.Helper()
	if len(a.Parsed) != want {
		t.Errorf("Parsed groups: got %d, want %d\n  Parsed: %v", len(a.Parsed), want, a.Parsed)
	}
}

// assertMatches checks that the returned matches equal the expected set
// (order-insensitive).
func assertMatches(t *testing.T, a vars.Analysis, want []string) {
	t.Helper()
	got := a.Matches
	if got == nil {
		got = []string{}
	}
	if want == nil {
		want = []string{}
	}
	sortedGot := make([]string, len(got))
	copy(sortedGot, got)
	sort.Strings(sortedGot)
	sortedWant := make([]string, len(want))
	copy(sortedWant, want)
	sort.Strings(sortedWant)
	if !reflect.DeepEqual(sortedGot, sortedWant) {
		t.Errorf("Matches mismatch:\n  got  (%d): %v\n  want (%d): %v", len(got), got, len(want), want)
	}
}

// assertMatchCount checks only the number of matches (for tests with many matches).
func assertMatchCount(t *testing.T, a vars.Analysis, want int) {
	t.Helper()
	if len(a.Matches) != want {
		t.Errorf("Match count: got %d, want %d", len(a.Matches), want)
	}
}

// preReleaseVersions is shared across per-language pre-release tests.
var preReleaseVersions = []string{
	"9.0.0-preview.1.24081.5",
	"9.0.0-preview.2.24128.5",
	"9.0.0-preview.3.24172.9",
	"9.0.0-rc.1.24452.1",
	"9.0.0-rc.2.24474.1",
	"9.0.0",
	"9.0.1",
	"9.1.0-preview.1",
	"10.0.0-preview.1",
	"10.0.0",
}

var wantPreRelease = []string{
	"9.0.0-preview.1.24081.5",
	"9.0.0-preview.2.24128.5",
	"9.0.0-preview.3.24172.9",
	"9.0.0-rc.1.24452.1",
}

// ─── Edge cases ──────────────────────────────────────────────────────────────

func TestEmptyConstraint(t *testing.T) {
	for _, style := range []vars.Style{vars.StyleNPM, vars.StylePy, vars.StyleNuGet, vars.StyleMaven, vars.StyleRuby, vars.StyleRust, vars.StyleGo} {
		t.Run(fmt.Sprintf("%s_empty", style), func(t *testing.T) {
			a := AnalyzeConstraint(style, "", vars.TestVersions)
			assertParsedCount(t, a, 0)
			assertMatches(t, a, []string{})
		})
	}
}

func TestUnknownStyle(t *testing.T) {
	a := AnalyzeConstraint("unknown", ">=1.0.0", vars.TestVersions)
	if a.Parsed != nil {
		t.Errorf("expected nil parsed for unknown style, got %v", a.Parsed)
	}
	assertMatches(t, a, []string{})
}
