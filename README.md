# versions

`versions` is a simple go module designed to manage version information for different ecosystem like pypi.org, maven.org,
npmjs.com, crates.io ect. in a structured manner.

## Motivation

There are a lot of tools in different ecosystem .e.g `github.com/go-simpler/goversion` and `github.com/Masterminds/semver`
in golang, `semver` in python and obviously there will be ecosystem specific package available to manage, parse and do a
lot of operations on version strings of any package.

But the real life scenario is different, when dealing with multiple ecosystem and encountered different versioning
of different packages, it becomes difficult to manage and parse them in a structured manner dut to non-semverish
versioning followed by different maintainers.

So, here comes `versions` to rescue, it provides a structured way to manage version information for different ecosystem
in a structured manner, do operations on them and parse them that cannot be parsed by these packages.

## Installation

```bash
  go get github.com/rng70/versions
```

## Usage
```go
package main

import (
    "fmt"
    "log"
    "github.com/rng70/versions/semver"
)

var version = "ver21.31.41beta23.alpha10.final-esm.21-27-gbf4dd2f-3-g04f7740-1-gff89e43-221-gff89e54+imcompatible-build.1234-20230815123045-9090-gdeadbeef",

parsedVersion, err := semver.NewVersion(version)
if err != nil {
    log.Fatal(err)
}  

fmt.Println(parsedVersion)
```

## Supported Ecosystem
- pypi.org
- maven.org
- npmjs.com
- crates.io
- golang.org
- rubygems.org
