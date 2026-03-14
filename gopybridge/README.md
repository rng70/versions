# semverish

`semverish` is a Python package that exposes the core functionality of the Go library [`github.com/rng70/versions`](https://github.com/rng70/versions) to Python — via a CGo-compiled shared library loaded at runtime through [CFFI](https://cffi.readthedocs.io/).

It handles version sorting and constraint resolution across multiple package ecosystems, including cases where versions do not follow strict semver (e.g. epoch prefixes, timestamps, commit hashes, pre-release labels with non-standard separators).

## How it works

```
Python (cffi)  →  libpyversions.so  (CGo)  →  github.com/rng70/versions  (Go)
```

The `gopybridge/semverish/bridge.go` file is compiled into a C shared library (`libpyversions.so`) using `go build -buildmode=c-shared`. The Python module loads this `.so` at import time via CFFI and calls the exported C functions directly.

## Installation

```bash
pip install semverish
```

Requires Python >= 3.8 and Linux (POSIX). The `.so` is bundled in the wheel.

**Dependency:** [`cffi >= 1.15.0`](https://pypi.org/project/cffi/) is installed automatically.

## Usage

### Sort versions

```python
import semverish

semverish.natural_sorted_versions(
    ['1.0.0', '2.0.0', '1.2.0', '1.0.0-alpha', '1.0.0-beta.1', '1.0.0.beta2', 'rel-1.3.3', '1.3.4'],
    descending=False
)
# ['1.0.0-alpha', '1.0.0-beta.1', '1.0.0.beta2', '1.0.0', '1.2.0', 'rel-1.3.3', '1.3.4', '2.0.0']
```

`descending=True` reverses the order (newest first).

By default (`safe_parse=True`), version strings that contain no numeric component — such as branch names or plain words — are silently excluded from the result:

```python
semverish.natural_sorted_versions(
    ['v1.0.0', 'reservation', 'refs/heads/main', '2.0.0'],
)
# ['v1.0.0', '2.0.0']  — 'reservation' and 'refs/heads/main' are dropped

semverish.natural_sorted_versions(
    ['v1.0.0', 'reservation', 'refs/heads/main', '2.0.0'],
    safe_parse=False
)
# all four strings are returned, named-only versions sorted to the front
```

### Resolve version constraints

```python
import semverish

semverish.analyze_constraints(
    language='npm',
    constraints='>1.1 <=2.9',
    versions=['1.1.1', '3.0', '2.9.9', '2.9.0', '1.9.0', '2.8.1', '1.0.0', '2.0.0', '1.2.0',
              '1.0.0-alpha', '1.0.0-beta.1', '1.0.0.beta2']
)
# ['1.1.1', '2.9.9', '2.9.0', '1.9.0', '2.8.1', '2.0.0', '1.2.0']

semverish.analyze_constraints(
    language='python',
    constraints='>1.1,<=4.9',
    versions=['1.1.1', '3.0', '2.9.9', '2.9.0', '1.9.0', '2.8.1', '1.0.0',
              'rel-2.0.0', '1.2.0', '1.0.0-alpha', '1.0.0-beta.1', 'rel_1.0.0.beta2']
)
```

## API

### `natural_sorted_versions(versions, descending=False, safe_parse=True) -> list[str]`

Sorts a list of version strings using the ecosystem-agnostic comparison engine from `github.com/rng70/versions`.

| Parameter | Type | Default | Description |
|---|---|---|---|
| `versions` | `list[str]` | — | Version strings to sort |
| `descending` | `bool` | `False` | Sort newest-first when `True` |
| `safe_parse` | `bool` | `True` | When `True`, strings with no numeric version core (e.g. `"reservation"`, `"refs/heads/main"`) are excluded from the result. Set to `False` to include them. |

Returns a sorted `list[str]`.

### `analyze_constraints(language, constraints, versions) -> list[str]`

Filters a list of versions to those that satisfy the given constraint string, parsed according to the specified ecosystem's syntax.

| Parameter | Type | Description |
|---|---|---|
| `language` | `str` | Ecosystem identifier (see table below) |
| `constraints` | `str` | Constraint expression in the ecosystem's native syntax |
| `versions` | `list[str]` | Version strings to filter |

Returns a `list[str]` of matching versions.

## Supported Ecosystems

| Ecosystem | `language` values | Constraint examples |
|---|---|---|
| PyPI | `python`, `py` | `>=1.0,<2.0`, `~=1.4`, `!=1.3.0` |
| NuGet | `nuget`, `csharp`, `dotnet`, `cs` | `[1.0,2.0)`, `(,1.0]`, `>=1.0.0` |
| npm | `npm`, `node`, `nodejs`, `javascript`, `js` | `^1.0.0`, `~1.2.3`, `>=1.0.0 <2.0.0`, `1.x` |
| Maven | `maven`, `java` | `[1.0,2.0)`, `[1.0.0]`, `>=1.0.0` |

## Building from Source

Requires Go (>= 1.21) and Python (>= 3.8) with `cffi` installed.

```bash
cd gopybridge

# Compile the Go shared library
make build

# Build the Python wheel
make wheel

# Full clean build
make all
```

The `Makefile` also provides targets for publishing:

```bash
make test-upload   # upload to test.pypi.org
make upload        # upload to pypi.org
```

## Why a separate README?

This file lives in `gopybridge/` and is used as the `long_description` in `setup.py`, so it appears on the [`semverish` PyPI page](https://pypi.org/project/semverish/). The root `README.md` documents the Go library itself, targeting Go developers. This file targets Python developers using the `semverish` pip package.

## License

MIT — see [LICENSE](../LICENSE).
