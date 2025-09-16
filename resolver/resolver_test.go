package resolver

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	versions := []string{
		"1.0.0", "1.2.3", "1.2.4", "1.3.0",
		"2.0.1", "2.3.1", "2.4.4", "2.5.2", "2.9.9", "3.0.0",
		"latest",
	}
	fmt.Println("=== VERSIONS ===\n", versions)

	npmTests := []string{
		"1.0.0 - 2.9999.9999",
		">=1.0.2 <2.1.2",
		">1.0.2 <=2.3.4",
		"2.0.1",
		"<1.0.0 || >=2.3.1 <2.4.5 || >=2.5.2 <3.0.0",
		"http://npmjs.com/example.tar.gz",
		"~1.2",
		"~1.2.3",
		"2.x",
		"3.3.x",
		"latest",
		"npm:pkg@1.0.0",
	}

	fmt.Println("=== NPM ===")
	for _, t := range npmTests {
		a := AnalyzeConstraint(StyleNPM, t, versions)
		fmt.Printf("\nRAW: %q\nPARSED: %v\nMATCHES: %v\n", a.Raw, a.Parsed, a.Matches)
	}

	pyTests := []string{
		">=1.0.2,<2.1.2",
		">1.0.2,<=2.3.4",
		"==2.0.1",
		"==1.2.*",
		"~=1.2.3",
		"~=1.4",
		"!=2.3.1,>=1.0.0,<3.0.0",
	}

	fmt.Println("\n=== Python (requirements.txt style) ===")
	for _, t := range pyTests {
		a := AnalyzeConstraint(StylePy, t, versions)
		fmt.Printf("\nRAW: %q\nPARSED: %v\nMATCHES: %v\n", a.Raw, a.Parsed, a.Matches)
	}

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
		a := AnalyzeConstraint(StyleNuGet, t, versions)
		fmt.Printf("\nRAW: %q\nPARSED: %v\nMATCHES: %v\n", a.Raw, a.Parsed, a.Matches)
	}

	// Example: sorted matches for complex NPM example
	ex := AnalyzeConstraint(StyleNPM, "<1.0.0 || >=2.3.1 <2.4.5 || >=2.5.2 <3.0.0", versions)
	// sort.Slice(ex.Matches, func(i, j int) bool { return cmpVersion(ex.Matches[i], ex.Matches[j]) < 0 })
	fmt.Println("\nSorted matches (NPM complex example):", ex.Matches)
}
