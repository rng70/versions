package canonicalized_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/rng70/versions/canonicalized"
	"github.com/rng70/versions/semver"
)

// ─── helpers ──────────────────────────────────────────────────────────────────

// GetGitTags runs: git ls-remote --tags --refs <repo>
// and returns tag names as []string.
// repo can be a local path (.) or a remote URL.
func GetGitTags(repo string) ([]string, error) {
	if strings.TrimSpace(repo) == "" {
		repo = "."
	}

	cmd := exec.Command("git", "ls-remote", "--tags", "--refs", repo)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return nil, fmt.Errorf("git ls-remote failed: %s", msg)
	}

	var tags []string
	for _, line := range strings.Split(stdout.String(), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Format: <hash>\trefs/tags/<tag>
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}
		ref := parts[1]
		tag := strings.TrimPrefix(ref, "refs/tags/")
		if tag != ref && tag != "" {
			tags = append(tags, tag)
		}
	}

	return tags, nil
}

// GetEcosystemsVersions fetches version numbers from ecosyste.ms.
// It handles both response shapes:
//   - {"versions": ["1.0.0", ...]}
//   - ["1.0.0", ...]
func GetEcosystemsVersions(ecosystem, pkg string, timeout time.Duration) ([]string, error) {
	ecosystem = strings.TrimSpace(ecosystem)
	pkg = strings.TrimSpace(pkg)
	if ecosystem == "" || pkg == "" {
		return nil, fmt.Errorf("ecosystem and package must be non-empty")
	}
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	escapedPkg := url.PathEscape(pkg)
	endpoint := fmt.Sprintf(
		"https://packages.ecosyste.ms/api/v1/registries/%s/packages/%s/version_numbers",
		ecosystem, escapedPkg,
	)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("ecosyste.ms error: HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	// Try dict form: {"versions":[...]}
	var dict struct {
		Versions []string `json:"versions"`
	}
	if err := json.Unmarshal(raw, &dict); err == nil && dict.Versions != nil {
		out := make([]string, 0, len(dict.Versions))
		for _, v := range dict.Versions {
			out = append(out, fmt.Sprint(v))
		}
		return out, nil
	}

	// Try list form: ["1.0.0", ...]
	var list []any
	if err := json.Unmarshal(raw, &list); err == nil {
		out := make([]string, 0, len(list))
		for _, v := range list {
			out = append(out, fmt.Sprint(v))
		}
		return out, nil
	}

	return nil, fmt.Errorf("unexpected ecosyste.ms response format: %s", string(raw))
}

// originalToCanonical builds a map[original]canonical for fast lookup.
func originalToCanonical(versions []*canonicalized.Version) map[string]string {
	m := make(map[string]string, len(versions))
	for _, v := range versions {
		m[v.Original] = v.Canonical
	}
	return m
}

// ─── tests ────────────────────────────────────────────────────────────────────

func TestGetGitTags(t *testing.T) {
	for _, tc := range versionMatchCases {
		t.Run(tc.name, func(t *testing.T) {
			tags, err := GetGitTags(tc.repo)
			if err != nil {
				t.Fatalf("GetGitTags error: %v", err)
			}
			if len(tags) == 0 {
				t.Error("expected at least one tag, got none")
			}
			t.Logf("fetched %d git tags", len(tags))
		})
	}
}

func TestGetEcosystemsVersions(t *testing.T) {
	for _, tc := range versionMatchCases {
		t.Run(tc.name, func(t *testing.T) {
			versions, err := GetEcosystemsVersions(tc.ecosystem, tc.pkg, testTimeout)
			if err != nil {
				t.Fatalf("GetEcosystemsVersions error: %v", err)
			}
			if len(versions) == 0 {
				t.Error("expected at least one version, got none")
			}
			t.Logf("fetched %d registry versions", len(versions))
		})
	}
}

func TestGetEcosystemsVersions_EmptyEcosystem(t *testing.T) {
	_, err := GetEcosystemsVersions("", "somepkg", testTimeout)
	if err == nil {
		t.Error("expected error for empty ecosystem, got nil")
	}
}

func TestGetEcosystemsVersions_EmptyPackage(t *testing.T) {
	_, err := GetEcosystemsVersions("pypi.org", "", testTimeout)
	if err == nil {
		t.Error("expected error for empty package, got nil")
	}
}

func TestGetGitTags_EmptyRepo(t *testing.T) {
	// Empty string defaults to "." — may succeed or fail depending on CWD.
	// We only assert the function does not panic.
	_, _ = GetGitTags("")
}

// TestVersionMatching is a table-driven test that, for each versionMatchCase:
//  1. Fetches live git tags and registry versions.
//  2. Canonicalises both lists into original→canonical maps.
//  3. For every pair {Registry: a, GitTag: b} in wantMatch:
//     - asserts a is present in the fetched registry versions
//     - asserts b is present in the fetched git tags
//     - asserts canonical(a) == canonical(b)
//
// The test fails for a pair if either original string is not found in its
// source, or if the two canonical forms differ.
func TestVersionMatching(t *testing.T) {
	for _, tc := range versionMatchCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rawTags, err := GetGitTags(tc.repo)
			if err != nil {
				t.Fatalf("GetGitTags: %v", err)
			}

			rawVersions, err := GetEcosystemsVersions(tc.ecosystem, tc.pkg, testTimeout)
			if err != nil {
				t.Fatalf("GetEcosystemsVersions: %v", err)
			}

			tagCanon := originalToCanonical(semver.NewVersionFromList(rawTags))
			regCanon := originalToCanonical(semver.NewVersionFromList(rawVersions))

			t.Logf("git tags: %d  |  registry versions: %d", len(tagCanon), len(regCanon))

			allPassed := true
			for _, pair := range tc.wantMatch {
				canonTag, inTags := tagCanon[pair.GitTag]
				canonReg, inReg := regCanon[pair.Registry]

				label := fmt.Sprintf("registry=%q git=%q", pair.Registry, pair.GitTag)

				switch {
				case !inReg && !inTags:
					t.Errorf("MISSING in both sources:  %s", label)
					allPassed = false
				case !inReg:
					t.Errorf("MISSING in registry:      %s  (git canonical=%q)", label, canonTag)
					allPassed = false
				case !inTags:
					t.Errorf("MISSING in git tags:      %s  (registry canonical=%q)", label, canonReg)
					allPassed = false
				case canonReg != canonTag:
					t.Errorf("CANONICAL MISMATCH:       %s  registry=%q  git=%q", label, canonReg, canonTag)
					allPassed = false
				default:
					t.Logf("MATCH  %s  canonical=%q", label, canonReg)
				}
			}

			if allPassed {
				t.Logf("all %d pairs matched", len(tc.wantMatch))
			}
		})
	}
}

// ─── vars ────────────────────────────────────────────────────────────────────
const testTimeout = 10 * time.Second

// versionPair is an expected match between a registry version string and a git
// tag string.  The test passes for this pair when both strings are found in
// their respective sources AND they canonicalize to the same value.
//
// Example: {Registry: "1.24.0", GitTag: "1.24.0"}
// Example: {Registry: "1.0",    GitTag: "v1.0"}   ← different originals, same canonical
type versionPair struct {
	Registry string // original version string as it appears on the registry
	GitTag   string // original tag string as it appears in the git repo
}

// versionMatchCase describes a parameterised version-matching test.
type versionMatchCase struct {
	name      string
	repo      string        // git repo URL passed to GetGitTags
	ecosystem string        // ecosyste.ms registry identifier
	pkg       string        // package name on the registry
	wantMatch []versionPair // pairs that must match between registry and git tags
}

// versionMatchCases is the single source of truth for all matching tests.
// Add a new entry here to test a different repo / registry combination.
var versionMatchCases = []versionMatchCase{
	{
		name:      "boto3/pypi",
		repo:      "https://github.com/boto/boto3",
		ecosystem: "pypi.org",
		pkg:       "boto3",
		wantMatch: []versionPair{
			{Registry: "1.24.0", GitTag: "1.24.0"},
			{Registry: "1.26.0", GitTag: "1.26.0"},
			{Registry: "1.28.0", GitTag: "1.28.0"},
			{Registry: "1.34.0", GitTag: "1.34.0"},
		},
	},
	{
		name:      "commons-lang/maven",
		repo:      "https://github.com/apache/commons-lang",
		ecosystem: "repo1.maven.org",
		pkg:       "org.apache.commons:commons-lang3",
		wantMatch: []versionPair{
			{Registry: "3.0", GitTag: "LANG_3_0"},
			{Registry: "3.0.1", GitTag: "LANG_3_0_1"},
			{Registry: "3.1", GitTag: "LANG_3_1"},
			{Registry: "3.10", GitTag: "rel/commons-lang-3.10"},
			{Registry: "3.11", GitTag: "rel/commons-lang-3.11"},
			{Registry: "3.12.0", GitTag: "rel/commons-lang-3.12.0"},
			{Registry: "3.13.0", GitTag: "rel/commons-lang-3.13.0"},
			{Registry: "3.14.0", GitTag: "rel/commons-lang-3.14.0"},
			{Registry: "3.15.0", GitTag: "rel/commons-lang-3.15.0"},
			{Registry: "3.16.0", GitTag: "rel/commons-lang-3.16.0"},
			{Registry: "3.17.0", GitTag: "rel/commons-lang-3.17.0"},
			{Registry: "3.18.0", GitTag: "rel/commons-lang-3.18.0"},
			{Registry: "3.19.0", GitTag: "rel/commons-lang-3.19.0"},
			{Registry: "3.2", GitTag: "LANG_3_2"},
			{Registry: "3.20.0", GitTag: "rel/commons-lang-3.20.0"},
			{Registry: "3.2.1", GitTag: "LANG_3_2_1"},
			{Registry: "3.3", GitTag: "LANG_3_3"},
			{Registry: "3.3.1", GitTag: "LANG_3_3_1"},
			{Registry: "3.3.2", GitTag: "LANG_3_3_2"},
			{Registry: "3.4", GitTag: "LANG_3_4"},
			{Registry: "3.5", GitTag: "LANG_3_5"},
			{Registry: "3.6", GitTag: "LANG_3_6"},
			{Registry: "3.7", GitTag: "LANG_3_7"},
			{Registry: "3.8", GitTag: "LANG_3_8"},
			{Registry: "3.8.1", GitTag: "LANG_3_8_1"},
			{Registry: "3.9", GitTag: "commons-lang-3.9"},
		},
	},
}
