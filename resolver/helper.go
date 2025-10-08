package resolver

import (
	"fmt"

	"github.com/rng70/versions/vars"
)

func NPMStyleTest() {
	npmTests := []string{
		"2.x",
		"3.3.x",
		"2.0.1",
		"~1.2",
		"~1.2.3",
		"1.0.0 - 2.9999.9999",
		">=1.0.2 <2.1.2",
		">1.0.2 <=2.3.4",
		"<1.0.0 || >=2.3.1 <2.4.5 || >=2.5.2 <3.0.0",
		"latest",
		"npm:pkg@1.0.0",
		"http://npmjs.com/example.tar.gz",
	}

	fmt.Println("=== NPM ===")
	for _, t := range npmTests {
		a := AnalyzeConstraint(vars.StyleNPM, t, vars.TestVersions)
		fmt.Printf("\nRAW: %q\nPARSED: %v\nMATCHES: %v\n", a.Raw, a.Parsed, a.Matches)
	}
}

func SpecialNPMStyleTest() {
	// Example: sorted matches for complex NPM example
	ex := AnalyzeConstraint(vars.StyleNPM, "<1.0.0 || >=2.3.1 <2.4.5 || >=2.5.2 <3.0.0", vars.TestVersions)
	// sort.Slice(ex.Matches, func(i, j int) bool { return cmpVersion(ex.Matches[i], ex.Matches[j]) < 0 })
	fmt.Println("\nSorted matches (NPM complex example):", ex.Matches)
}

func PythonStyleTestHelper() {
	pyTests := []string{
		"~=1.4",
		"==1.2.*",
		"==2.0.1",
		"~=1.2.3",
		">=1.0.2,<2.1.2",
		">1.0.2,<=2.3.4",
		"!=2.3.1,>=1.0.0,<3.0.0",
		"2.3.1",
	}

	fmt.Println("\n=== Python (requirements.txt style) ===")
	for _, t := range pyTests {
		a := AnalyzeConstraint(vars.StylePy, t, vars.TestVersions)
		fmt.Printf("\nRAW: %q\nPARSED: %v\nMATCHES: %v\n", a.Raw, a.Parsed, a.Matches)
	}
}

func NuGetStyleTest() {
	nugetTests := []string{
		"[1.0.0,2.0.0)",
		"(,2.4.5]",
		"[2.5.2,)",
		"[2.0.1]",
		"1.*",
		"1.2.*",
		"2.3.1",
	}

	fmt.Println("\n=== NuGet (C#) ===")
	for _, t := range nugetTests {
		a := AnalyzeConstraint(vars.StyleNuGet, t, vars.TestVersions)
		fmt.Printf("\nRAW: %q\nPARSED: %v\nMATCHES: %v\n", a.Raw, a.Parsed, a.Matches)
	}
}
