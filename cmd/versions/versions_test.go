package main

import (
	"testing"

	"github.com/rng70/versions/canonicalized"
	"github.com/rng70/versions/semver"
	"github.com/rng70/versions/utils"
)

func Test(t *testing.T) {
	versions := utils.SingleTestList

	out := semver.NewVersionFromList(versions)

	// write to a file
	_error := utils.WriteToFile("testdata/output/parsed.json", out)
	if _error != nil {
		panic(_error)
	}

	_error = utils.WriteToFileWithMinimalContext("testdata/output/minimal:parsed.json", out)
	if _error != nil {
		panic(_error)
	}

	// sort descending
	canonicalized.SortVersions(out, false)

	_error = utils.WriteToFile("testdata/output/sorted:parsed.json", out)
	if _error != nil {
		panic(_error)
	}

	_error = utils.WriteToFileWithMinimalContext("testdata/output/sorted:minimal:parsed.json", out)
	if _error != nil {
		panic(_error)
	}
}
