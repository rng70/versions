import pytest
from semverish import natural_sorted_versions, analyze_constraints

# ─── Dataset ──────────────────────────────────────────────────────────────────

GENERIC_SEMVER_VERSIONS = [
    "3.0.0",
    "3.0.0-rc.2",
    "3.0.0-rc.1",
    "3.0.0-beta.1",
    "3.0.0-alpha.1",
    "2.0.0",
    "1.0.0",
    "1.1.0",
    "2.0.0-preview.1",
]

# Rules: alpha < beta < preview < rc < stable within same core version;
# core versions ordered numerically across groups.
GENERIC_SEMVER_VERSIONS_SORTED_ASCENDING = [
    "1.0.0",
    "1.1.0",
    "2.0.0-preview.1",
    "2.0.0",
    "3.0.0-alpha.1",
    "3.0.0-beta.1",
    "3.0.0-rc.1",
    "3.0.0-rc.2",
    "3.0.0",
]

GENERIC_SEMVER_VERSIONS_SORTED_DESCENDING = [
    "3.0.0",
    "3.0.0-rc.2",
    "3.0.0-rc.1",
    "3.0.0-beta.1",
    "3.0.0-alpha.1",
    "2.0.0",
    "2.0.0-preview.1",
    "1.1.0",
    "1.0.0",
]


# ─── natural_sorted_versions (mirrors SortedVersions) ─────────────────────────

class TestNaturalSortedVersions:

    def test_ascending_default(self):
        result = natural_sorted_versions(["3.0.0", "1.0.0", "2.0.0"])
        assert len(result) == 3
        assert result[0] == "1.0.0"
        assert result[2] == "3.0.0"

    def test_explicit_ascending(self):
        result = natural_sorted_versions(["3.0.0", "1.0.0", "2.0.0"], descending=False)
        assert result[0] == "1.0.0"

    def test_descending(self):
        result = natural_sorted_versions(["1.0.0", "3.0.0", "2.0.0"], descending=True)
        assert len(result) == 3
        assert result[0] == "3.0.0"
        assert result[2] == "1.0.0"

    def test_empty(self):
        result = natural_sorted_versions([])
        assert len(result) == 0

    def test_single(self):
        result = natural_sorted_versions(["2.5.0"])
        assert len(result) == 1
        assert result[0] == "2.5.0"

    def test_preserves_original_strings(self):
        # natural_sorted_versions must return the original strings, not canonical ones
        input_versions = ["v3.0.0", "v1.0.0", "v2.0.0"]
        result = natural_sorted_versions(input_versions, descending=False)
        for v in result:
            assert v in input_versions, f"{v!r} was not found in original inputs"

    def test_v_prefix_preserved(self):
        result = natural_sorted_versions(["v3.0.0", "v1.0.0", "v2.0.0"], descending=False)
        for v in result:
            assert v.startswith("v"), f"v-prefix not preserved: {v!r}"

    def test_stable_after_prerelease_ascending(self):
        result = natural_sorted_versions(["1.0.0", "1.0.0-beta1"], descending=False)
        assert result[0] == "1.0.0-beta1"
        assert result[1] == "1.0.0"

    def test_prerelease_order_ascending(self):
        input_versions = ["1.0.0", "1.0.0-rc1", "1.0.0-beta1", "1.0.0-alpha1"]
        result = natural_sorted_versions(input_versions, descending=False)
        # ascending: alpha < beta < rc < stable
        assert result[-1] == "1.0.0", f"last in ascending should be stable, got {result[-1]!r}"
        assert "alpha" in result[0], f"first in ascending should be alpha, got {result[0]!r}"

    def test_major_version_order(self):
        input_versions = ["10.0.0", "9.0.0", "11.0.0", "8.0.0"]
        result = natural_sorted_versions(input_versions, descending=False)
        assert result == ["8.0.0", "9.0.0", "10.0.0", "11.0.0"]

    def test_safe_parse_default_true_filters_named(self):
        input_versions = ["v1.0.0", "reservation", "refs/heads/smallish-refactor", "2.0.0"]
        result = natural_sorted_versions(input_versions)  # safe_parse=True by default
        assert "reservation" not in result
        assert "refs/heads/smallish-refactor" not in result
        assert len(result) == 2

    def test_safe_parse_false_keeps_named(self):
        input_versions = ["v1.0.0", "reservation", "refs/heads/smallish-refactor", "2.0.0"]
        result = natural_sorted_versions(input_versions, descending=False, safe_parse=False)
        assert len(result) == 4

    def test_safe_parse_true_keeps_named(self):
            input_versions = ["v1.0.0", "reservation", "refs/heads/smallish-refactor", "2.0.0"]
            result = natural_sorted_versions(input_versions, descending=True, safe_parse=True)
            assert len(result) == 2

    def test_mixed_formats(self):
        input_versions = ["v1.0.0", "2.0.0-rc1", "3.0.0", "latest"]
        # safe_parse=False to keep all entries including "latest"
        result = natural_sorted_versions(input_versions, safe_parse=False)
        assert len(result) == 4

    def test_large_set_ascending(self):
        input_versions = [
            "9.0.0", "9.0.0-preview.1.24081.5", "9.0.0-rc.1.24452.1",
            "8.0.0", "8.0.0-preview.1.23112.2",
            "10.0.0", "10.0.0-preview.1.25120.3",
        ]
        result = natural_sorted_versions(input_versions, descending=False)
        assert len(result) == len(input_versions)
        # last should be 10.0.0 (highest stable)
        assert result[-1] == "10.0.0", f"last in ascending should be 10.0.0, got {result[-1]!r}"
        # first should be a prerelease
        assert any(c in result[0] for c in ["-", "preview", "rc"]), (
            f"first in ascending large set should not be stable: {result[0]!r}"
        )

    def test_invalid_input_raises(self):
        with pytest.raises(RuntimeError, match="ERROR"):
            natural_sorted_versions("not-a-list")  # type: ignore[arg-type]

    def test_returns_list(self):
        result = natural_sorted_versions(["2.0.0", "1.0.0"])
        assert isinstance(result, list)


# ─── Predefined order dataset (mirrors sorted_assert_test.go) ─────────────────

class TestSortedGenericPredefinedOrder:

    def test_ascending(self):
        got = natural_sorted_versions(GENERIC_SEMVER_VERSIONS, descending=False)
        want = GENERIC_SEMVER_VERSIONS_SORTED_ASCENDING
        assert len(got) == len(want), (
            f"ascending: length mismatch: got {len(got)}, want {len(want)}\n"
            f"got:  {got}\nwant: {want}"
        )
        for i, (g, w) in enumerate(zip(got, want)):
            assert g == w, f"ascending pos {i}: got {g!r}, want {w!r}"

    def test_descending(self):
        got = natural_sorted_versions(GENERIC_SEMVER_VERSIONS, descending=True)
        want = GENERIC_SEMVER_VERSIONS_SORTED_DESCENDING
        assert len(got) == len(want), (
            f"descending: length mismatch: got {len(got)}, want {len(want)}\n"
            f"got:  {got}\nwant: {want}"
        )
        for i, (g, w) in enumerate(zip(got, want)):
            assert g == w, f"descending pos {i}: got {g!r}, want {w!r}"

    def test_stable_last_in_ascending(self):
        got = natural_sorted_versions(GENERIC_SEMVER_VERSIONS, descending=False)
        last = got[-1]
        assert "-" not in last, f"last version in ascending should be stable, got {last!r}"

    def test_stable_first_in_descending(self):
        got = natural_sorted_versions(GENERIC_SEMVER_VERSIONS, descending=True)
        first = got[0]
        assert "-" not in first, f"first version in descending should be stable, got {first!r}"

    def test_prerelease_order_alpha_before_beta(self):
        got = natural_sorted_versions(GENERIC_SEMVER_VERSIONS, descending=False)
        alpha_pos = got.index("3.0.0-alpha.1")
        beta_pos = got.index("3.0.0-beta.1")
        assert alpha_pos < beta_pos, (
            f"alpha should come before beta: alpha at {alpha_pos}, beta at {beta_pos}"
        )

    def test_prerelease_order_rc_before_stable(self):
        got = natural_sorted_versions(GENERIC_SEMVER_VERSIONS, descending=False)
        rc_pos = got.index("3.0.0-rc.2")
        stable_pos = got.index("3.0.0")
        assert rc_pos < stable_pos, (
            f"rc should come before stable: rc at {rc_pos}, stable at {stable_pos}"
        )

    def test_rc1_before_rc2(self):
        got = natural_sorted_versions(GENERIC_SEMVER_VERSIONS, descending=False)
        rc1_pos = got.index("3.0.0-rc.1")
        rc2_pos = got.index("3.0.0-rc.2")
        assert rc1_pos < rc2_pos, (
            f"rc.1 should come before rc.2: rc.1 at {rc1_pos}, rc.2 at {rc2_pos}"
        )

    def test_cross_major_prerelease_after_lower_stable(self):
        # 1.1.0 (stable) must come before 2.0.0-preview.1 (pre-release of higher major)
        got = natural_sorted_versions(GENERIC_SEMVER_VERSIONS, descending=False)
        v110_pos = got.index("1.1.0")
        v200pre_pos = got.index("2.0.0-preview.1")
        assert v110_pos < v200pre_pos, (
            f"1.1.0 should come before 2.0.0-preview.1: "
            f"1.1.0 at {v110_pos}, 2.0.0-preview.1 at {v200pre_pos}"
        )


# ─── analyze_constraints ──────────────────────────────────────────────────────

class TestAnalyzeConstraints:

    def test_python_gte(self):
        result = analyze_constraints("python", ">=2.0.0", ["1.0.0", "2.0.0", "3.0.0"])
        assert "1.0.0" not in result, "1.0.0 should not match >=2.0.0"

    def test_npm_gte(self):
        result = analyze_constraints("npm", ">=2.0.0", ["1.0.0", "2.0.0", "3.0.0"])
        assert "1.0.0" not in result, "1.0.0 should not match >=2.0.0"

    def test_nuget_range(self):
        result = analyze_constraints("nuget", "[2.0.0, 3.0.0)", ["1.0.0", "2.0.0", "3.0.0"])
        assert len(result) > 0, "expected at least one match for [2.0.0, 3.0.0)"
        assert "1.0.0" not in result, "1.0.0 should not match [2.0.0, 3.0.0)"
        assert "3.0.0" not in result, "3.0.0 should not match [2.0.0, 3.0.0) (exclusive upper bound)"

    def test_language_aliases(self):
        versions = ["1.0.0", "2.0.0"]
        constraint = ">=1.0.0"
        aliases = ["python", "py", "npm", "node", "nodejs", "javascript", "js", "nuget", "csharp", "dotnet", "cs"]
        for lang in aliases:
            result = analyze_constraints(lang, constraint, versions)
            assert isinstance(result, list), f"language {lang!r} did not return a list"

    def test_unknown_language_defaults_to_python(self):
        result = analyze_constraints("unknown", ">=1.0.0", ["1.0.0", "2.0.0"])
        assert isinstance(result, list), "unknown language should fall back to python and return a list"

    def test_invalid_versions_raises(self):
        with pytest.raises(RuntimeError, match="ERROR"):
            analyze_constraints("python", ">=1.0.0", "not-a-list")  # type: ignore[arg-type]

    def test_empty_versions_list(self):
        result = analyze_constraints("python", ">=1.0.0", [])
        # Go marshals a nil slice as JSON null; Python decodes that as None
        assert not result, f"expected no matches for empty versions list, got {result}"

    def test_returns_list(self):
        result = analyze_constraints("python", ">=1.0.0", ["1.0.0", "2.0.0", "3.0.0"])
        assert isinstance(result, list)

    def test_python_range_constraint(self):
        result = analyze_constraints("python", ">1.1.0,<4.9.9", ["1.0.0", "2.0.0", "3.0.0", "5.0.0"])
        assert "1.0.0" not in result
        assert "5.0.0" not in result
        assert "2.0.0" in result
        assert "3.0.0" in result

    def test_npm_range_constraint(self):
        # npm hyphen range "a - b" is the AND form; ">=a <=b" is parsed as OR by this parser
        result = analyze_constraints("npm", "1.2.0 - 2.9.0", ["1.1.0", "3.0.0", "2.9.0", "1.2.0"])
        assert "1.1.0" not in result
        assert "3.0.0" not in result
        assert "1.2.0" in result
        assert "2.9.0" in result

    def test_py_alias_same_as_python(self):
        versions = ["1.0.0", "2.0.0", "3.0.0"]
        constraint = ">=2.0.0"
        assert analyze_constraints("py", constraint, versions) == \
               analyze_constraints("python", constraint, versions)

    def test_node_alias_same_as_npm(self):
        versions = ["1.0.0", "2.0.0", "3.0.0"]
        constraint = ">=2.0.0"
        assert analyze_constraints("node", constraint, versions) == \
               analyze_constraints("npm", constraint, versions)
