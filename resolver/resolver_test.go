package resolver

import (
	"fmt"
	"testing"

	"github.com/rng70/versions/vars"
)

func Test(t *testing.T) {

	fmt.Println("=== VERSIONS ===")
	fmt.Println(vars.TestVersions)

	// ************ Python Test ************ //
	PythonStyleTestHelper()

	// ************ NPM Test ************ //
	NPMStyleTest()
	SpecialNPMStyleTest()

	// ************ NuGet Test ************ //
	NuGetStyleTest()
}
