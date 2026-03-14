package semver

import (
	"testing"
)

// ─── Dataset ──────────────────────────────────────────────────────────────────

var GenericSemverVersions = []string{
	"3.0.0",
	"3.0.0-rc.2",
	"3.0.0-rc.1",
	"3.0.0-beta.1",
	"3.0.0-alpha.1",
	"2.0.0",
	"1.0.0",
	"1.1.0",
	"2.0.0-preview.1",
}

// Predefined expected ascending order for GenericSemverVersions.
// Rules verified: alpha < beta < preview < rc < stable within same core version;
// core versions ordered numerically across groups.
var GenericSemverVersionsSortedAscending = []string{
	"1.0.0",
	"1.1.0",
	"2.0.0-preview.1",
	"2.0.0",
	"3.0.0-alpha.1",
	"3.0.0-beta.1",
	"3.0.0-rc.1",
	"3.0.0-rc.2",
	"3.0.0",
}

var GenericSemverVersionsSortedDescending = []string{
	"3.0.0",
	"3.0.0-rc.2",
	"3.0.0-rc.1",
	"3.0.0-beta.1",
	"3.0.0-alpha.1",
	"2.0.0",
	"2.0.0-preview.1",
	"1.1.0",
	"1.0.0",
}

// ─── TestSorted_Generic_PredefinedOrder ───────────────────────────────────────

func TestSorted_Generic_PredefinedOrder(t *testing.T) {
	t.Run("ascending", func(t *testing.T) {
		got := SortedVersions(GenericSemverVersions, false)
		want := GenericSemverVersionsSortedAscending

		if len(got) != len(want) {
			t.Fatalf("ascending: length mismatch: got %d, want %d\ngot:  %v\nwant: %v", len(got), len(want), got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("ascending pos %d: got %q, want %q", i, got[i], want[i])
			}
		}
	})

	t.Run("descending", func(t *testing.T) {
		got := SortedVersions(GenericSemverVersions, true)
		want := GenericSemverVersionsSortedDescending

		if len(got) != len(want) {
			t.Fatalf("descending: length mismatch: got %d, want %d\ngot:  %v\nwant: %v", len(got), len(want), got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("descending pos %d: got %q, want %q", i, got[i], want[i])
			}
		}
	})

	t.Run("stable_last_in_ascending", func(t *testing.T) {
		got := SortedVersions(GenericSemverVersions, false)
		last := NewVersion(got[len(got)-1])
		if !last.IsStable() {
			t.Errorf("last version in ascending should be stable, got %q", got[len(got)-1])
		}
	})

	t.Run("stable_first_in_descending", func(t *testing.T) {
		got := SortedVersions(GenericSemverVersions, true)
		first := NewVersion(got[0])
		if !first.IsStable() {
			t.Errorf("first version in descending should be stable, got %q", got[0])
		}
	})

	t.Run("prerelease_order_alpha_before_beta", func(t *testing.T) {
		got := SortedVersions(GenericSemverVersions, false)
		// find positions of 3.0.0-alpha.1 and 3.0.0-beta.1
		alphaPos, betaPos := -1, -1
		for i, v := range got {
			switch v {
			case "3.0.0-alpha.1":
				alphaPos = i
			case "3.0.0-beta.1":
				betaPos = i
			}
		}
		if alphaPos == -1 || betaPos == -1 {
			t.Fatal("could not find alpha or beta in result")
		}
		if alphaPos >= betaPos {
			t.Errorf("alpha should come before beta: alpha at %d, beta at %d", alphaPos, betaPos)
		}
	})

	t.Run("prerelease_order_rc_before_stable", func(t *testing.T) {
		got := SortedVersions(GenericSemverVersions, false)
		rcPos, stablePos := -1, -1
		for i, v := range got {
			switch v {
			case "3.0.0-rc.2":
				rcPos = i
			case "3.0.0":
				stablePos = i
			}
		}
		if rcPos == -1 || stablePos == -1 {
			t.Fatal("could not find rc or stable 3.0.0 in result")
		}
		if rcPos >= stablePos {
			t.Errorf("rc should come before stable: rc at %d, stable at %d", rcPos, stablePos)
		}
	})

	t.Run("rc1_before_rc2", func(t *testing.T) {
		got := SortedVersions(GenericSemverVersions, false)
		rc1Pos, rc2Pos := -1, -1
		for i, v := range got {
			switch v {
			case "3.0.0-rc.1":
				rc1Pos = i
			case "3.0.0-rc.2":
				rc2Pos = i
			}
		}
		if rc1Pos == -1 || rc2Pos == -1 {
			t.Fatal("could not find rc.1 or rc.2 in result")
		}
		if rc1Pos >= rc2Pos {
			t.Errorf("rc.1 should come before rc.2: rc.1 at %d, rc.2 at %d", rc1Pos, rc2Pos)
		}
	})

	t.Run("cross_major_prerelease_after_lower_stable", func(t *testing.T) {
		got := SortedVersions(GenericSemverVersions, false)
		// 1.1.0 (stable) must come before 2.0.0-preview.1 (pre-release of higher major)
		v110Pos, v200prePos := -1, -1
		for i, v := range got {
			switch v {
			case "1.1.0":
				v110Pos = i
			case "2.0.0-preview.1":
				v200prePos = i
			}
		}
		if v110Pos == -1 || v200prePos == -1 {
			t.Fatal("could not find 1.1.0 or 2.0.0-preview.1 in result")
		}
		if v110Pos >= v200prePos {
			t.Errorf("1.1.0 should come before 2.0.0-preview.1: 1.1.0 at %d, 2.0.0-preview.1 at %d", v110Pos, v200prePos)
		}
	})
}
