package semver

import (
	"testing"
)

// ─── NewVersion ───────────────────────────────────────────────────────────────

func TestNewVersion_Simple(t *testing.T) {
	v := NewVersion("1.2.3")
	if v.Original != "1.2.3" {
		t.Errorf("original: got %q, want %q", v.Original, "1.2.3")
	}
	if v.Canonical == "" {
		t.Error("canonical should not be empty for a valid version")
	}
}

func TestNewVersion_VPrefix(t *testing.T) {
	v := NewVersion("v2.0.0")
	if v.Prefix != "v" {
		t.Errorf("prefix: got %q, want %q", v.Prefix, "v")
	}
	if v.Original != "v2.0.0" {
		t.Errorf("original: got %q, want %q", v.Original, "v2.0.0")
	}
}

func TestNewVersion_Prerelease(t *testing.T) {
	v := NewVersion("1.0.0-beta1")
	if v.IsStable() {
		t.Error("beta version should not be stable")
	}
}

func TestNewVersion_NoCore(t *testing.T) {
	v := NewVersion("latest")
	if v.Major == nil || *v.Major != -1 {
		t.Errorf("no-core: expected major=-1, got %v", v.Major)
	}
}

func TestNewVersion_PreservesOriginalUnchanged(t *testing.T) {
	input := "v3.0.0-preview.1.24081.5"
	v := NewVersion(input)
	if v.Original != input {
		t.Errorf("original: got %q, want %q", v.Original, input)
	}
}

// ─── NewVersionFromList ───────────────────────────────────────────────────────

func TestNewVersionFromList_Empty(t *testing.T) {
	vs := NewVersionFromList([]string{})
	if len(vs) != 0 {
		t.Errorf("empty list: expected 0 versions, got %d", len(vs))
	}
}

func TestNewVersionFromList_Single(t *testing.T) {
	vs := NewVersionFromList([]string{"1.0.0"})
	if len(vs) != 1 {
		t.Fatalf("expected 1 version, got %d", len(vs))
	}
	if vs[0].Original != "1.0.0" {
		t.Errorf("original: got %q, want %q", vs[0].Original, "1.0.0")
	}
}

func TestNewVersionFromList_Multiple(t *testing.T) {
	input := []string{"1.0.0", "2.0.0", "3.0.0-beta1"}
	vs := NewVersionFromList(input)
	if len(vs) != 3 {
		t.Fatalf("expected 3 versions, got %d", len(vs))
	}
	for i, v := range vs {
		if v.Original != input[i] {
			t.Errorf("pos %d original: got %q, want %q", i, v.Original, input[i])
		}
	}
}

func TestNewVersionFromList_ReturnsPointers(t *testing.T) {
	vs := NewVersionFromList([]string{"1.0.0", "2.0.0"})
	for i, v := range vs {
		if v == nil {
			t.Errorf("pos %d: pointer should not be nil", i)
		}
	}
}

func TestNewVersionFromList_MixedFormats(t *testing.T) {
	input := []string{"v1.0.0", "2.0.0-rc1", "3.0.0", "latest"}
	vs := NewVersionFromList(input)
	if len(vs) != 4 {
		t.Fatalf("expected 4 versions, got %d", len(vs))
	}
}

// ─── SortedParsedVersions ─────────────────────────────────────────────────────

func TestSortedParsedVersions_AscendingDefault(t *testing.T) {
	vs := SortedParsedVersions([]string{"3.0.0", "1.0.0", "2.0.0"})
	if len(vs) != 3 {
		t.Fatalf("expected 3 versions, got %d", len(vs))
	}
	if vs[0].Original != "1.0.0" {
		t.Errorf("first: got %q, want %q", vs[0].Original, "1.0.0")
	}
	if vs[2].Original != "3.0.0" {
		t.Errorf("last: got %q, want %q", vs[2].Original, "3.0.0")
	}
}

func TestSortedParsedVersions_ExplicitAscending(t *testing.T) {
	vs := SortedParsedVersions([]string{"3.0.0", "1.0.0", "2.0.0"}, false)
	if vs[0].Original != "1.0.0" {
		t.Errorf("first: got %q, want %q", vs[0].Original, "1.0.0")
	}
}

func TestSortedParsedVersions_Descending(t *testing.T) {
	vs := SortedParsedVersions([]string{"1.0.0", "3.0.0", "2.0.0"}, true)
	if len(vs) != 3 {
		t.Fatalf("expected 3 versions, got %d", len(vs))
	}
	if vs[0].Original != "3.0.0" {
		t.Errorf("first (desc): got %q, want %q", vs[0].Original, "3.0.0")
	}
	if vs[2].Original != "1.0.0" {
		t.Errorf("last (desc): got %q, want %q", vs[2].Original, "1.0.0")
	}
}

func TestSortedParsedVersions_Empty(t *testing.T) {
	vs := SortedParsedVersions([]string{})
	if len(vs) != 0 {
		t.Errorf("empty: expected 0 versions, got %d", len(vs))
	}
}

func TestSortedParsedVersions_Single(t *testing.T) {
	vs := SortedParsedVersions([]string{"1.0.0"})
	if len(vs) != 1 || vs[0].Original != "1.0.0" {
		t.Errorf("single: unexpected result %v", vs)
	}
}

func TestSortedParsedVersions_PreReleaseOrder(t *testing.T) {
	input := []string{"1.0.0", "1.0.0-rc1", "1.0.0-beta1", "1.0.0-alpha1"}
	vs := SortedParsedVersions(input, false) // ascending
	// ascending: alpha < beta < rc < stable
	if !vs[len(vs)-1].IsStable() {
		t.Errorf("last in ascending should be stable, got %q", vs[len(vs)-1].Original)
	}
	if !vs[0].IsAlpha() {
		t.Errorf("first in ascending should be alpha, got %q", vs[0].Original)
	}
}

func TestSortedParsedVersions_MajorVersionOrder(t *testing.T) {
	input := []string{"10.0.0", "9.0.0", "11.0.0", "8.0.0"}
	vs := SortedParsedVersions(input, false)
	expected := []string{"8.0.0", "9.0.0", "10.0.0", "11.0.0"}
	for i, e := range expected {
		if vs[i].Original != e {
			t.Errorf("pos %d: got %q, want %q", i, vs[i].Original, e)
		}
	}
}

// ─── SortedVersions ───────────────────────────────────────────────────────────

func TestSortedVersions_AscendingDefault(t *testing.T) {
	result := SortedVersions([]string{"3.0.0", "1.0.0", "2.0.0"})
	if len(result) != 3 {
		t.Fatalf("expected 3 versions, got %d", len(result))
	}
	if result[0] != "1.0.0" {
		t.Errorf("first: got %q, want %q", result[0], "1.0.0")
	}
	if result[2] != "3.0.0" {
		t.Errorf("last: got %q, want %q", result[2], "3.0.0")
	}
}

func TestSortedVersions_Descending(t *testing.T) {
	result := SortedVersions([]string{"1.0.0", "3.0.0", "2.0.0"}, true)
	if len(result) != 3 {
		t.Fatalf("expected 3 versions, got %d", len(result))
	}
	if result[0] != "3.0.0" {
		t.Errorf("first (desc): got %q, want %q", result[0], "3.0.0")
	}
	if result[2] != "1.0.0" {
		t.Errorf("last (desc): got %q, want %q", result[2], "1.0.0")
	}
}

func TestSortedVersions_Empty(t *testing.T) {
	result := SortedVersions([]string{})
	if len(result) != 0 {
		t.Errorf("empty: expected 0 versions, got %d", len(result))
	}
}

func TestSortedVersions_Single(t *testing.T) {
	result := SortedVersions([]string{"2.5.0"})
	if len(result) != 1 || result[0] != "2.5.0" {
		t.Errorf("single: unexpected result %v", result)
	}
}

func TestSortedVersions_PreservesOriginalStrings(t *testing.T) {
	// SortedVersions must return the Original strings, not canonical ones
	input := []string{"v3.0.0", "v1.0.0", "v2.0.0"}
	result := SortedVersions(input, false)
	for _, r := range result {
		found := false
		for _, o := range input {
			if r == o {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("result %q was not found in original inputs %v", r, input)
		}
	}
}

func TestSortedVersions_StableAfterPrerelease(t *testing.T) {
	input := []string{"1.0.0", "1.0.0-beta1"}
	result := SortedVersions(input, false)
	if result[0] != "1.0.0-beta1" || result[1] != "1.0.0" {
		t.Errorf("pre-release should come before stable: got %v", result)
	}
}

func TestSortedVersions_SafeParse_DefaultTrue_FiltersNamed(t *testing.T) {
	input := []string{
		"v1.0.0",
		"reservation",
		"refs/heads/smallish-refactor",
		"2.0.0",
	}
	result := SortedVersions(input) // safeParse=true by default
	for _, r := range result {
		if r == "reservation" || r == "refs/heads/smallish-refactor" {
			t.Errorf("safeParse=true should have excluded %q but it was returned", r)
		}
	}
	if len(result) != 2 {
		t.Errorf("expected 2 versions after filtering named-only, got %d: %v", len(result), result)
	}
}

func TestSortedVersions_SafeParse_False_KeepsNamed(t *testing.T) {
	input := []string{
		"v1.0.0",
		"reservation",
		"refs/heads/smallish-refactor",
		"2.0.0",
	}
	result := SortedVersions(input, false, false) // desc=false, safeParse=false
	if len(result) != 4 {
		t.Errorf("safeParse=false should keep all 4 entries, got %d: %v", len(result), result)
	}
}

func TestSortedParsedVersions_SafeParse_DefaultTrue_FiltersNamed(t *testing.T) {
	input := []string{"1.0.0", "reservation", "refs/heads/main", "2.0.0"}
	vs := SortedParsedVersions(input)
	if len(vs) != 2 {
		t.Errorf("expected 2 parsed versions, got %d", len(vs))
	}
	for _, v := range vs {
		if v.Major == nil || *v.Major == -1 {
			t.Errorf("safeParse=true should have excluded named-only version %q", v.Original)
		}
	}
}

func TestSortedVersions_LargeSet(t *testing.T) {
	input := []string{
		"9.0.0", "9.0.0-preview.1.24081.5", "9.0.0-rc.1.24452.1",
		"8.0.0", "8.0.0-preview.1.23112.2",
		"10.0.0", "10.0.0-preview.1.25120.3",
	}
	result := SortedVersions(input, false) // ascending
	if len(result) != len(input) {
		t.Fatalf("length mismatch: got %d, want %d", len(result), len(input))
	}
	// first should be a preview/prerelease
	first := NewVersion(result[0])
	if first.IsStable() {
		t.Errorf("first in ascending large set should not be stable: %q", result[0])
	}
	// last should be 10.0.0 (highest stable)
	if result[len(result)-1] != "10.0.0" {
		t.Errorf("last in ascending should be 10.0.0, got %q", result[len(result)-1])
	}
}
