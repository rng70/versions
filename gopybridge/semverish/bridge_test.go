package main

import (
	"encoding/json"
	"strings"
	"testing"
)

// ─── naturalSortedVersions ────────────────────────────────────────────────────

func TestNaturalSortedVersions_AscendingOrder(t *testing.T) {
	input := `["3.0.0","1.0.0","2.0.0"]`
	result := naturalSortedVersions(input, false)

	var got []string
	if err := json.Unmarshal([]byte(result), &got); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 versions, got %d", len(got))
	}
	if got[0] != "1.0.0" || got[2] != "3.0.0" {
		t.Errorf("ascending order: got %v", got)
	}
}

func TestNaturalSortedVersions_DescendingOrder(t *testing.T) {
	input := `["1.0.0","3.0.0","2.0.0"]`
	result := naturalSortedVersions(input, true)

	var got []string
	if err := json.Unmarshal([]byte(result), &got); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 versions, got %d", len(got))
	}
	if got[0] != "3.0.0" || got[2] != "1.0.0" {
		t.Errorf("descending order: got %v", got)
	}
}

func TestNaturalSortedVersions_EmptyArray(t *testing.T) {
	result := naturalSortedVersions(`[]`, false)
	var got []string
	if err := json.Unmarshal([]byte(result), &got); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty array, got %v", got)
	}
}

func TestNaturalSortedVersions_SingleVersion(t *testing.T) {
	result := naturalSortedVersions(`["1.2.3"]`, false)
	var got []string
	if err := json.Unmarshal([]byte(result), &got); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if len(got) != 1 || got[0] != "1.2.3" {
		t.Errorf("single version: got %v", got)
	}
}

func TestNaturalSortedVersions_PreReleaseOrder(t *testing.T) {
	input := `["1.0.0","1.0.0-beta1","1.0.0-alpha1"]`
	result := naturalSortedVersions(input, false)

	var got []string
	if err := json.Unmarshal([]byte(result), &got); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	// ascending: alpha < beta < stable
	if got[len(got)-1] != "1.0.0" {
		t.Errorf("last in ascending should be stable 1.0.0, got %q", got[len(got)-1])
	}
}

func TestNaturalSortedVersions_InvalidJSON(t *testing.T) {
	result := naturalSortedVersions(`not-json`, false)
	if !strings.HasPrefix(result, "ERROR: invalid json:") {
		t.Errorf("expected ERROR prefix, got %q", result)
	}
}

func TestNaturalSortedVersions_InvalidJSONObject(t *testing.T) {
	result := naturalSortedVersions(`{"key":"value"}`, false)
	if !strings.HasPrefix(result, "ERROR: invalid json:") {
		t.Errorf("expected ERROR prefix for object input, got %q", result)
	}
}

func TestNaturalSortedVersions_PreservesOriginalStrings(t *testing.T) {
	input := `["v3.0.0","v1.0.0","v2.0.0"]`
	result := naturalSortedVersions(input, false)

	var got []string
	if err := json.Unmarshal([]byte(result), &got); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	for _, v := range got {
		if !strings.HasPrefix(v, "v") {
			t.Errorf("v-prefix not preserved in result: %q", v)
		}
	}
}

func TestNaturalSortedVersions_ReturnsValidJSON(t *testing.T) {
	input := `["2.0.0","1.0.0"]`
	result := naturalSortedVersions(input, false)
	if !json.Valid([]byte(result)) {
		t.Errorf("result is not valid JSON: %q", result)
	}
}

// ─── languageStyle ────────────────────────────────────────────────────────────

func TestLanguageStyle_Python(t *testing.T) {
	for _, lang := range []string{"python", "py"} {
		s := languageStyle(lang)
		if s != "python" {
			t.Errorf("languageStyle(%q): got %q, want %q", lang, s, "python")
		}
	}
}

func TestLanguageStyle_NuGet(t *testing.T) {
	for _, lang := range []string{"nuget", "csharp", "dotnet", "cs"} {
		s := languageStyle(lang)
		if s != "nuget" {
			t.Errorf("languageStyle(%q): got %q, want %q", lang, s, "nuget")
		}
	}
}

func TestLanguageStyle_NPM(t *testing.T) {
	for _, lang := range []string{"npm", "node", "nodejs", "javascript", "js"} {
		s := languageStyle(lang)
		if s != "npm" {
			t.Errorf("languageStyle(%q): got %q, want %q", lang, s, "npm")
		}
	}
}

func TestLanguageStyle_Maven(t *testing.T) {
	for _, lang := range []string{"maven", "java"} {
		s := languageStyle(lang)
		// maven and java map to StyleNPM in the bridge
		if s != "npm" {
			t.Errorf("languageStyle(%q): got %q, want %q", lang, s, "npm")
		}
	}
}

func TestLanguageStyle_UnknownDefaultsToPython(t *testing.T) {
	for _, lang := range []string{"", "ruby", "rust", "unknown"} {
		s := languageStyle(lang)
		if s != "python" {
			t.Errorf("languageStyle(%q): got %q, want default %q", lang, s, "python")
		}
	}
}

// ─── analyzeConstraints ───────────────────────────────────────────────────────

func TestAnalyzeConstraints_PythonGte(t *testing.T) {
	versions := `["1.0.0","2.0.0","3.0.0"]`
	result := analyzeConstraints("python", ">=2.0.0", versions)

	var got []string
	if err := json.Unmarshal([]byte(result), &got); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	for _, v := range got {
		if v == "1.0.0" {
			t.Errorf("1.0.0 should not match >=2.0.0, but it did")
		}
	}
}

func TestAnalyzeConstraints_NPM(t *testing.T) {
	versions := `["1.0.0","2.0.0","3.0.0"]`
	result := analyzeConstraints("npm", ">=2.0.0", versions)

	var got []string
	if err := json.Unmarshal([]byte(result), &got); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	for _, v := range got {
		if v == "1.0.0" {
			t.Errorf("1.0.0 should not match >=2.0.0, but it did")
		}
	}
}

func TestAnalyzeConstraints_NuGet(t *testing.T) {
	versions := `["1.0.0","2.0.0","3.0.0"]`
	result := analyzeConstraints("nuget", "[2.0.0, 3.0.0)", versions)

	var got []string
	if err := json.Unmarshal([]byte(result), &got); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if len(got) == 0 {
		t.Error("expected at least one match for [2.0.0, 3.0.0)")
	}
	for _, v := range got {
		if v == "1.0.0" || v == "3.0.0" {
			t.Errorf("%q should not match [2.0.0, 3.0.0)", v)
		}
	}
}

func TestAnalyzeConstraints_LanguageAliases(t *testing.T) {
	versions := `["1.0.0","2.0.0"]`
	constraint := ">=1.0.0"
	aliases := []string{"python", "py", "npm", "node", "nodejs", "javascript", "js", "nuget", "csharp", "dotnet", "cs"}
	for _, lang := range aliases {
		result := analyzeConstraints(lang, constraint, versions)
		if strings.HasPrefix(result, "ERROR:") {
			t.Errorf("languageStyle(%q) returned error: %q", lang, result)
		}
		if !json.Valid([]byte(result)) {
			t.Errorf("languageStyle(%q) returned invalid JSON: %q", lang, result)
		}
	}
}

func TestAnalyzeConstraints_UnknownLanguageDefaultsToPython(t *testing.T) {
	versions := `["1.0.0","2.0.0"]`
	result := analyzeConstraints("unknown", ">=1.0.0", versions)
	if strings.HasPrefix(result, "ERROR:") {
		t.Errorf("unknown language should fall back to python, got error: %q", result)
	}
	if !json.Valid([]byte(result)) {
		t.Errorf("unknown language result is not valid JSON: %q", result)
	}
}

func TestAnalyzeConstraints_InvalidVersionsJSON(t *testing.T) {
	result := analyzeConstraints("python", ">=1.0.0", `not-json`)
	if !strings.HasPrefix(result, "ERROR: invalid json:") {
		t.Errorf("expected ERROR prefix, got %q", result)
	}
}

func TestAnalyzeConstraints_EmptyVersionsList(t *testing.T) {
	result := analyzeConstraints("python", ">=1.0.0", `[]`)
	var got []string
	if err := json.Unmarshal([]byte(result), &got); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected no matches for empty versions list, got %v", got)
	}
}

func TestAnalyzeConstraints_ReturnsValidJSON(t *testing.T) {
	versions := `["1.0.0","2.0.0","3.0.0"]`
	result := analyzeConstraints("python", ">=1.0.0", versions)
	if !json.Valid([]byte(result)) {
		t.Errorf("result is not valid JSON: %q", result)
	}
}
