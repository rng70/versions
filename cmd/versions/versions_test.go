package main

import (
	"testing"

	"github.com/rng70/versions/canonicalized"
	"github.com/rng70/versions/semver"
	"github.com/rng70/versions/utils"
)

func Test(t *testing.T) {
	versions := utils.SingleTestList

	parsed := semver.NewVersionFromList(versions)
	sortedParsed := semver.SortedParsedVersions(versions, true)

	// write to a file
	_error := utils.WriteToFile("testdata/output/parsed.json", parsed)
	if _error != nil {
		panic(_error)
	}

	_error = utils.WriteToFileWithMinimalContext("testdata/output/minimal:parsed.json", sortedParsed)
	if _error != nil {
		panic(_error)
	}

	// sort descending
	canonicalized.SortVersions(parsed)
	sortedMinimalParsed := semver.SortedVersions(versions, true)

	_error = utils.WriteToFile("testdata/output/sorted:parsed.json", parsed)
	if _error != nil {
		panic(_error)
	}

	_error = utils.WriteToFile("testdata/output/sorted:minimal:parsed.json", sortedMinimalParsed)
	if _error != nil {
		panic(_error)
	}
}
