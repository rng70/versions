package main

import (
	"fmt"

	"github.com/rng70/versions/resolver"
	"github.com/rng70/versions/vars"
)

func main() {
	fmt.Println("=== VERSIONS ===")
	fmt.Println(vars.TestVersions)

	// ************ Python Test ************ //
	resolver.PythonStyleTestHelper()

	// ************ NPM Test ************ //
	resolver.NPMStyleTest()
	resolver.SpecialNPMStyleTest()

	// ************ NuGet Test ************ //
	resolver.NuGetStyleTest()

	// *********** Random Test ************ //
	resolver.RandomTest()
}
