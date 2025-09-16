package main

import (
	"sort"
	"testing"

	"github.com/rng70/versions/canonicalized"
	"github.com/rng70/versions/semver"
	"github.com/rng70/versions/utils"
)

func Test(t *testing.T) {
	versions := utils.FullList

	out := semver.NewVersionFromList(versions)

	// write to a file
	_error := utils.WriteToFile("testdata/output/out_test_parsed.json", out)
	if _error != nil {
		panic(_error)
	}

	// flatten all parsed versions
	all := make([]*canonicalized.Version, 0)
	for _, item := range out {
		all = append(all, &item.Parsed)
	}

	// sort ascending (default)
	canonicalized.SortVersions(all)

	// if you want to re-order `out` slice itself:
	sort.Slice(out, func(i, j int) bool {
		return out[i].Parsed.LessThan(&out[j].Parsed)
	})

	_error = utils.WriteToFile("testdata/output/out_sorted-test_parsed.json", out)
	if _error != nil {
		panic(_error)
	}

	_error = utils.WriteToFileWithMinimalContext("testdata/output/out_test_sorted-parsed-minimal.json", out)
	if _error != nil {
		panic(_error)
	}
}
