# Test Results & Coverage — `canonicalized`

> Run command (excludes live-network tests):
> ```bash
> go test -v -coverprofile=cover.out ./canonicalized/... \
>   -skip "^(TestGetGitTags$|TestGetEcosystemsVersions$|TestVersionMatching$)"
> go tool cover -func=cover.out
> ```

## Summary

| Metric | Value |
|---|---|
| Total tests | 125 |
| Passed | 125 |
| Failed | 0 |
| Statement coverage | **82.3 %** |

---

## Test suites

### `ParseVersionString` (23 tests)

| Test | What it checks |
|---|---|
| `TestParseVersionString_ThreePart` | Standard `1.2.3` → correct Major/Minor/Patch, canonical, original |
| `TestParseVersionString_OriginalPreserved` | Original string is stored verbatim |
| `TestParseVersionString_VPrefix` | `v2.0.1` → prefix `"v"`, Revision=0, canonical `2.0.1` |
| `TestParseVersionString_TwoPart` | `4.7` → Patch defaults to 0, canonical `4.7.0` |
| `TestParseVersionString_OnePart` | `5` → Minor/Patch default to 0, canonical `5.0.0` |
| `TestParseVersionString_FourPart_Revision` | `1.2.3.4` → Revision=4, canonical `1.2.3-4` |
| `TestParseVersionString_EmptyString` | `""` → all fields -1, canonical `""` |
| `TestParseVersionString_NoVersionCore` | `"latest"` → stored as Prefix, Major=-1, canonical `""` |
| `TestParseVersionString_BetaWithNumericTag` | `1.0.0-beta1` → Type `beta`, Tag=1, canonical `1.0.0-beta.1` |
| `TestParseVersionString_AlphaWithNumericTag` | `2.0.0-alpha2` → Type `alpha`, Tag=2, canonical `2.0.0-alpha.2` |
| `TestParseVersionString_RCWithNumericTag` | `3.0.0-rc1` → Type `rc`, Tag=1, canonical `3.0.0-rc.1` |
| `TestParseVersionString_PreviewDotNotation` | `4.0.0-preview.1` → Type `preview`, Tag=1 |
| `TestParseVersionString_ComplexPreview` | `9.0.0-preview.1.24081.5` → intermediate numerics preserved in canonical |
| `TestParseVersionString_ComplexRC` | `9.0.0-rc.2.24474.1.12234` → all prerelease tokens in canonical |
| `TestParseVersionString_BuildMetadata` | `1.0.0+build.123` → metadata preserved, canonical `1.0.0+build.123` |
| `TestParseVersionString_MavenFinalSuffix` | `4.1.0.Final` → canonical `4.1.0-final` |
| `TestParseVersionString_MavenFinalSuffixWithNumber` | `4.1.0.Final.418` → canonical `4.1.0-final.418` |
| `TestParseVersionString_StableHasNoType` | `2.0.0` → Type slice is empty |
| `TestParseVersionString_AliasA_IsAlpha` | `1.0.0-a1` → alias `a` maps to `alpha`, canonical `1.0.0-alpha.1` |
| `TestParseVersionString_AliasB_IsBeta` | `1.0.0-b2` → alias `b` maps to `beta`, canonical `1.0.0-beta.2` |
| `TestVersion_Lang3_0` | `COMMON_LANG_3_0` → underscore core `3_0` → Major=3 Minor=0 Patch=0 |
| `TestVersion_Lang3_0_1` | `COMMON_LANG_3_0_1` → underscore core `3_0_1` → Patch=1, canonical `3.0.1` |
| `TestVersion_Lang3_0_1_0` | `COMMON_LANG_3_0_1_0` → 4-part underscore, Revision=0, canonical `3.0.1` |
| `TestVersion_Lang3_0_1_2` | `COMMON_LANG_3_0_1_2` → 4-part underscore, Revision=2, canonical `3.0.1-2` |

### `NewVersion` (5 tests)

| Test | What it checks |
|---|---|
| `TestNewVersion_Simple` | `1.2.3` parses correctly via `NewVersion` |
| `TestNewVersion_VPrefix` | `v3.0.0` → Prefix `"v"` preserved |
| `TestNewVersion_Prerelease` | `1.0.0-rc1` → `IsRC()` returns true |
| `TestNewVersion_NoCore` | `"latest"` → Major=-1 |
| `TestNewVersion_PreservesOriginal` | `v2.0.0-beta1` → Original stored verbatim |

### `Compare` (13 tests)

| Test | What it checks |
|---|---|
| `TestCompare_EqualVersions` | Same version → 0 |
| `TestCompare_MajorLess` / `MajorGreater` | Major ordering |
| `TestCompare_MinorLess` / `MinorGreater` | Minor ordering |
| `TestCompare_PatchLess` / `PatchGreater` | Patch ordering |
| `TestCompare_StableGreaterThanPrerelease` | `1.0.0` > `1.0.0-beta1` |
| `TestCompare_PrereleaseEqual` | Same pre-release → 0 |
| `TestCompare_AlphaBeforeBeta` | alpha < beta |
| `TestCompare_BetaBeforeRC` | beta < rc |
| `TestCompare_PreviewBeforeRC` | preview < rc |
| `TestCompare_RCBeforeStable` | rc < stable |
| `TestCompare_TwoPreviewNumbers` | `preview.1` < `preview.2` |
| `TestCompare_TwoRCNumbers` | `rc.1` < `rc.2` |
| `TestCompare_DifferentMajorSamePrerelease` | Major dominates pre-release |

### `LessThan` / `GreaterThan` / `Equal` / `LessThanOrEqual` / `GreaterThanOrEqual` (19 tests)

Each operator is verified for the lower, equal, and greater cases, plus a stable-vs-prerelease edge case where applicable.

### `Prerelease` (6 tests)

Verifies `Prerelease()` returns `""` for stable and a non-empty string for alpha, beta, rc, preview, and versions with an `Extra` field set.

### `MetadataStr` / `String` (11 tests)

Covers empty, single, and multi-entry `Metadata` slices, and `String()` output for versions built from components, with prefix, with revision, with prerelease tags, and with build metadata.

### `CompareType` (7 tests)

Equal types, different name, different tag, different length, both stable, different `Extra`, one `Extra` nil.

### `IsStable` / `IsAlpha` / `IsBeta` / `IsRC` / `IsPreview` / `IsPseudo` (22 tests)

Each predicate is tested for a positive match and for common negative cases (e.g., `IsBeta` on a stable, alpha, and RC version).

### `SortVersions` (8 tests)

Ascending, descending, empty slice, single element, pre-release ordering within same core, mixed major versions, explicit `false`/no-arg (both ascending).

### Integration — input validation (3 tests)

| Test | What it checks |
|---|---|
| `TestGetEcosystemsVersions_EmptyEcosystem` | Returns error for empty ecosystem |
| `TestGetEcosystemsVersions_EmptyPackage` | Returns error for empty package name |
| `TestGetGitTags_EmptyRepo` | Defaults to `.` without panicking |

---

## Coverage by function

| File | Function | Coverage |
|---|---|---|
| `compare.go` | `NewVersion` | 100.0 % |
| `compare.go` | `newVersion` | 100.0 % |
| `compares.go` | `Compare` | 94.1 % |
| `compares.go` | `LessThan` | 100.0 % |
| `compares.go` | `GreaterThan` | 100.0 % |
| `compares.go` | `Equal` | 100.0 % |
| `compares.go` | `LessThanOrEqual` | 100.0 % |
| `compares.go` | `GreaterThanOrEqual` | 100.0 % |
| `compares.go` | `Prerelease` | 100.0 % |
| `compares.go` | `MetadataStr` | 100.0 % |
| `compares.go` | `String` | 100.0 % |
| `compares.go` | `CompareType` | 100.0 % |
| `compares.go` | `IsStable` | 100.0 % |
| `compares.go` | `IsAlpha` | 100.0 % |
| `compares.go` | `IsBeta` | 100.0 % |
| `compares.go` | `IsRC` | 100.0 % |
| `compares.go` | `IsPreview` | 100.0 % |
| `compares.go` | `IsPseudo` | 100.0 % |
| `compares.go` | `isPrerelease` | 100.0 % |
| `compares.go` | `safeInt` | 66.7 % |
| `compares.go` | `cmpPtrInt` | 100.0 % |
| `compares.go` | `stageWeight` | 71.4 % |
| `compares.go` | `isNum` | 100.0 % |
| `compares.go` | `compareId` | 71.4 % |
| `compares.go` | `comparePrerelease` | 75.0 % |
| `compares.go` | `toPreIdents` | 100.0 % |
| `compares.go` | `i64` | 100.0 % |
| `compares.go` | `SortVersions` | 90.0 % |
| `parser.go` | `ParseVersionString` | 82.7 % |
| `utils.go` | `trimWeirdQuotes` | 100.0 % |
| `utils.go` | `removeNoise` | 100.0 % |
| `utils.go` | `parseCoreInts` | 72.7 % |
| `utils.go` | `isNumeric` | 0.0 % |
| `utils.go` | `splitAlphaNumTail` | 90.0 % |
| `utils.go` | `canonicalName` | 75.0 % |
| `utils.go` | `findCoreIndex` | 95.5 % |
| `utils.go` | `takeAlphaPrefix` | 100.0 % |
| `utils.go` | `isOnlyAlpha` | 66.7 % |
| `utils.go` | `parseTimestampToISO` | 0.0 % |
| `utils.go` | `splitTokens` | 100.0 % |
| `utils.go` | `extractCommitsAndTags` | 30.0 % |
| `utils.go` | `extractCommitsGOnly` | 33.3 % |
| `utils.go` | `setMajorMinorPatch` | 85.7 % |
| **Total** | | **82.3 %** |

### Notable gaps

| Function | Coverage | Reason |
|---|---|---|
| `parseTimestampToISO` | 0 % | No test exercises a version string with a 8- or 14-digit timestamp token |
| `isNumeric` | 0 % | Internal helper reached only through paths not yet exercised by the current suite |
| `extractCommitsAndTags` / `extractCommitsGOnly` | 30–33 % | Git-describe version strings (e.g. `1.2.3-4-gabcdef`) are not yet in the unit tests |

---

## Live-network tests (excluded above)

The following tests require outbound network access and are excluded from the standard coverage run. Run them explicitly with:

```bash
go test -v ./canonicalized/... -run "^(TestGetGitTags$|TestGetEcosystemsVersions$|TestVersionMatching$)"
```

| Test | What it checks |
|---|---|
| `TestGetGitTags/boto3/pypi` | Fetches real git tags from `github.com/boto/boto3` |
| `TestGetGitTags/commons-lang/maven` | Fetches real git tags from `github.com/apache/commons-lang` |
| `TestGetEcosystemsVersions/boto3/pypi` | Fetches boto3 versions from `packages.ecosyste.ms` |
| `TestGetEcosystemsVersions/commons-lang/maven` | Fetches commons-lang3 versions from `packages.ecosyste.ms` |
| `TestVersionMatching/boto3/pypi` | Asserts 4 known `(registry, git-tag)` pairs canonicalize identically |
| `TestVersionMatching/commons-lang/maven` | Asserts 26 known pairs across dot-notation and `LANG_X_Y` tag formats |
