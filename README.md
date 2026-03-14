# versions

`versions` is a Go library for parsing, comparing, sorting, and resolving version constraints across multiple package ecosystems — npm, PyPI, NuGet, Maven, RubyGems, crates.io, and Go modules.

## Motivation

Every package ecosystem invents its own versioning conventions. Tools like `Masterminds/semver` or Python's `packaging` work well within a single ecosystem, but break down when you need to handle multiple ecosystems uniformly — especially when packages use non-standard, non-semver version strings (epoch prefixes, git describe output, timestamps, commit hashes, etc.).

`versions` solves this by providing:
- A unified version parser that handles exotic version strings across all ecosystems
- Ecosystem-aware constraint parsers (`^1.0`, `~>2.0`, `[1.0,2.0)`, `~=1.4`, etc.)
- A single resolver API that maps constraints to matching versions regardless of ecosystem

## Installation

```bash
go get github.com/rng70/versions
```

## Packages

| Package | Description |
|---|---|
| `canonicalized` | Core version struct and comparison engine |
| `parser` | Ecosystem-specific constraint parsers |
| `resolver` | Constraint resolution (parses + filters) |
| `semver` | Version list utilities (parse, sort) |
| `vars` | Shared types (`Constraint`, `Analysis`, `Style`) |

## Usage

### Parse a version string

```go
import "github.com/rng70/versions/semver"

v := semver.NewVersion("ver21.31.41beta23.alpha10.final-esm.21-27-gbf4dd2f+build.1234")
fmt.Println(v.Original)  // original input
fmt.Println(v.Canonical) // normalized form
fmt.Println(v.IsStable()) // false
```

### Sort a list of versions

```go
import "github.com/rng70/versions/semver"

versions := []string{"2.0.0", "1.0.0-alpha", "1.0.0", "3.0.0-rc.1"}

// Returns sorted original strings (descending), safeParse=true by default
sorted := semver.SortedVersions(versions, true)
// ["3.0.0-rc.1", "2.0.0", "1.0.0", "1.0.0-alpha"]

// Returns sorted parsed Version structs
parsed := semver.SortedParsedVersions(versions, true)
```

Both functions accept an optional second variadic bool for `safeParse` (default `true`). When enabled, version strings that contain no numeric component — such as branch refs or plain words — are excluded from the result:

```go
mixed := []string{"v1.0.0", "reservation", "refs/heads/main", "2.0.0"}

// safeParse=true (default): named-only strings dropped
semver.SortedVersions(mixed)               // ["v1.0.0", "2.0.0"]
semver.SortedVersions(mixed, false)        // same, ascending

// safeParse=false: all strings kept
semver.SortedVersions(mixed, false, false) // ["reservation", "refs/heads/main", "v1.0.0", "2.0.0"]
```

### Resolve version constraints

```go
import (
    "github.com/rng70/versions/resolver"
    "github.com/rng70/versions/vars"
)

available := []string{"1.0.0", "1.5.0", "2.0.0", "3.0.0-beta.1"}

// npm
result := resolver.AnalyzeConstraint(vars.StyleNPM, "^1.0.0", available)
fmt.Println(result.Matches) // ["1.0.0", "1.5.0"]

// Python (PEP-440)
result = resolver.AnalyzeConstraint(vars.StylePy, ">=1.0.0,<2.0.0", available)
fmt.Println(result.Matches) // ["1.0.0", "1.5.0"]

// NuGet
result = resolver.AnalyzeConstraint(vars.StyleNuGet, "[1.0.0, 2.0.0)", available)
fmt.Println(result.Matches) // ["1.0.0", "1.5.0"]

// Maven
result = resolver.AnalyzeConstraint(vars.StyleMaven, "[1.0.0,2.0.0)", available)
fmt.Println(result.Matches) // ["1.0.0", "1.5.0"]

// Ruby
result = resolver.AnalyzeConstraint(vars.StyleRuby, "~> 1.0", available)
fmt.Println(result.Matches) // ["1.0.0", "1.5.0"]

// Rust (Cargo)
result = resolver.AnalyzeConstraint(vars.StyleRust, ">=1.0.0, <2.0.0", available)
fmt.Println(result.Matches) // ["1.0.0", "1.5.0"]

// Go modules
result = resolver.AnalyzeConstraint(vars.StyleGo, ">=v1.0.0, <v2.0.0", available)
fmt.Println(result.Matches) // ["1.0.0", "1.5.0"]
```

The returned `vars.Analysis` contains:

```go
type Analysis struct {
    Raw     string          // original constraint string
    Parsed  [][]Constraint  // parsed constraint groups (OR of ANDs)
    Matches []string        // versions that satisfy the constraint
}
```

### Parse constraints directly

```go
import "github.com/rng70/versions/parser"

// Returns [][]vars.Constraint: outer = OR groups, inner = AND constraints
groups, err := parser.ParseNPM("^1.2.3 || >=2.0.0 <3.0.0")

// Filter a version list manually
matches := parser.FilterMatches(groups, available)
```

## Supported Ecosystems

| Ecosystem | Style constant | Constraint examples |
|---|---|---|
| npmjs.com | `vars.StyleNPM` | `^1.0.0`, `~1.2.3`, `>=1.0.0 <2.0.0`, `1.x`, `*` |
| pypi.org | `vars.StylePy` | `>=1.0,<2.0`, `~=1.4`, `==1.2.*`, `!=1.3.0` |
| nuget.org | `vars.StyleNuGet` | `[1.0,2.0)`, `(,1.0]`, `1.0.*`, `>=1.0.0` |
| maven.org | `vars.StyleMaven` | `[1.0,2.0)`, `[1.0.0]`, `>=1.0.0` |
| rubygems.org | `vars.StyleRuby` | `~> 2.0`, `~> 2.0.3`, `>= 1.0.0` |
| crates.io | `vars.StyleRust` | `^1.0.0`, `~1.2.3`, `>=1.0.0, <2.0.0`, `1.*` |
| golang.org | `vars.StyleGo` | `>=v1.0.0`, `>=v1.0.0, <v2.0.0` |

## Version struct

The `canonicalized.Version` struct exposes all parsed components:

```go
type Version struct {
    Original   string           // original input string
    Canonical  string           // normalized semver-ish form
    Prefix     string           // e.g. "v", "ver"
    Major      *int64
    Minor      *int64
    Patch      *int64
    Revision   *int64           // 4th numeric component
    Type       []TypeTag        // pre-release stages (alpha, beta, rc, ...)
    Metadata   []BuildMetadata  // build metadata after "+"
    Timestamp  []TimestampInfo  // parsed timestamps (YYYYMMDD or YYYYMMDDhhmmss)
    CommitHash []CommitHashInfo // git commit hashes
}
```

Predicate methods: `IsStable()`, `IsAlpha()`, `IsBeta()`, `IsRC()`, `IsPreview()`, `IsPseudo()`

Comparison methods: `Compare()`, `Equal()`, `LessThan()`, `GreaterThan()`, `LessThanOrEqual()`, `GreaterThanOrEqual()`
