package vars

import "regexp"

type Constraint struct {
	Op  string
	Ver string
}

type Analysis struct {
	Raw     string
	Parsed  [][]Constraint
	Matches []string
}

type Style string

const (
	StyleNPM   Style = "npm"
	StylePy    Style = "python"
	StyleNuGet Style = "nuget"
	StyleMaven Style = "maven"
	StyleRuby  Style = "ruby"
	StyleRust  Style = "rust"
	StyleGo    Style = "go"
)

type ConstraintResult struct {
	Raw     string
	Parsed  string
	Matches []string
}

// Notes: Go's regexp does not support non-capturing groups (?:...), so parentheses are used.
// We capture two main groups: operator (group 1) and token (group 2).
var (
	ReDashRange = regexp.MustCompile(`^\s*([0-9]+(?:\.[0-9]+){2})\s*-\s*([0-9]+(?:\.[0-9]+){2})\s*$`)
	//ReNpmToken   = regexp.MustCompile(`(?i)(<=|>=|<|>|=|~|\^)?\s*([0-9]+(\.[0-9]+){0,2}(\.(x|X|\*))?|latest|npm:[^\s@]+@\d+(\.\d+){0,2}|https?://\S+|file:[^\s]+|\*|[0-9]+)`)
	ReNpmToken = regexp.MustCompile(
		`(?i)(?:` +
			`(<=|>=|<|>|=|~|\^)?\s*` +
			`([0-9]+(\.[0-9]+)*(\.[A-Za-z][A-Za-z0-9]*)?(-[A-Za-z0-9]+(\.[A-Za-z0-9]+)*)?|latest|npm:[^\s@]+@\d+(\.\d+){0,2}|https?://\S+|file:[^\s]+|\*|[0-9]+)` +
			`)(?:\s+` +
			`(<=|>=|<|>|=|~|\^)\s*` +
			`([0-9]+(\.[0-9]+)*(\.[A-Za-z][A-Za-z0-9]*)?(-[A-Za-z0-9]+(\.[A-Za-z0-9]+)*)?)` +
			`)?`,
	)
	//RePyPart     = regexp.MustCompile(`^(==|!=|<=|>=|<|>|~=|===)?\s*([0-9]+(\.[0-9]+){0,2}(\.\*)?)\s*$`)
	RePyPart = regexp.MustCompile(
		`^(==|!=|<=|>=|<|>|~=|===)?\s*` +
			`(` +
			`[0-9]+` +
			`(?:\.[0-9]*){0,2}` +
			`(?:` +
			`\.[A-Za-z][A-Za-z0-9]*` +
			`|[A-Za-z][A-Za-z0-9]*` +
			`)?` +
			`(?:\.[A-Za-z][A-Za-z0-9]*)*` +
			`(?:-[A-Za-z0-9]+(?:\.[A-Za-z0-9]+)*)?` +
			`)` +
			`(?:\+[^\s]+)?` +
			`\s*$`,
	)

	ReNuGetRange = regexp.MustCompile(`^\s*([\[\(])\s*([^,\s]*)\s*,\s*([^\]\)\s]*)\s*([\]\)])\s*$`)
)

var (
	//TestVersions = []string{
	//	"1.0.0",
	//	"1.2.3",
	//	"1.2.4",
	//	"1.3.0",
	//	"2.0.1",
	//	"2.3.1",
	//	"2.4.4",
	//	"2.5.2",
	//	"2.9.9",
	//	"3.0.0",
	//	"latest",
	//}

	TestVersions = []string{
		"10.0.1",
		"10.0.1-beta1",
		"10.0.2",
		"10.0.3",
		"11.0.1",
		"11.0.1-beta1",
		"11.0.1-beta2",
		"11.0.1-beta3",
		"11.0.2",
		"12.0.1",
		"12.0.1-beta1",
		"12.0.1-beta2",
		"12.0.2",
		"12.0.2-beta1",
		"12.0.2-beta2",
		"12.0.2-beta3",
		"12.0.3",
		"12.0.3-beta1",
		"12.0.3-beta2",
		"13.0.1",
		"13.0.1-beta1",
		"13.0.1-beta2",
		"13.0.2",
		"13.0.2-beta1",
		"13.0.2-beta2",
		"13.0.2-beta3",
		"13.0.3",
		"13.0.3-beta1",
		"13.0.4",
		"13.0.4-beta1",
		"13.0.5-beta1",
		"3.5.8",
		"4.0.1",
		"4.0.2",
		"4.0.3",
		"4.0.4",
		"4.0.5",
		"4.0.6",
		"4.0.7",
		"4.0.8",
		"4.5.1",
		"4.5.10",
		"4.5.11",
		"4.5.2",
		"4.5.3",
		"4.5.4",
		"4.5.5",
		"4.5.6",
		"4.5.7",
		"4.5.8",
		"4.5.9",
		"5.0.1",
		"5.0.2",
		"5.0.3",
		"5.0.4",
		"5.0.5",
		"5.0.6",
		"5.0.7",
		"5.0.8",
		"6.0.1",
		"6.0.1-beta1",
		"6.0.2",
		"6.0.3",
		"6.0.4",
		"6.0.5",
		"6.0.6",
		"6.0.7",
		"6.0.8",
		"7.0.1",
		"7.0.1-beta1",
		"7.0.1-beta2",
		"7.0.1-beta3",
		"8.0.1",
		"8.0.1-beta1",
		"8.0.1-beta2",
		"8.0.1-beta3",
		"8.0.1-beta4",
		"8.0.2",
		"8.0.3",
		"8.0.4-beta1",
		"9.0.1",
		"9.0.1-beta1",
		"9.0.2-beta1",
		"9.0.2-beta2",
	}

	TestVersionsNettyCodecHttp2 = []string{
		"4.1.0.Beta4",
		"4.1.0.Beta5",
		"4.1.0.Beta6",
		"4.1.0.Beta7",
		"4.1.0.Beta8",
		"4.1.0.CR1",
		"4.1.0.CR2",
		"4.1.0.CR3",
		"4.1.0.CR4",
		"4.1.0.CR5",
		"4.1.0.CR6",
		"4.1.0.CR7",
		"4.1.0.Final",
		"4.1.100.Final",
		"4.1.101.Final",
		"4.1.102.Final",
		"4.1.103.Final",
		"4.1.104.Final",
		"4.1.105.Final",
		"4.1.106.Final",
		"4.1.107.Final",
		"4.1.108.Final",
		"4.1.109.Final",
		"4.1.10.Final",
		"4.1.110.Final",
		"4.1.111.Final",
		"4.1.112.Final",
		"4.1.113.Final",
		"4.1.114.Final",
		"4.1.115.Final",
		"4.1.116.Final",
		"4.1.117.Final",
		"4.1.118.Final",
		"4.1.119.Final",
		"4.1.11.Final",
		"4.1.120.Final",
		"4.1.121.Final",
		"4.1.122.Final",
		"4.1.123.Final",
		"4.1.124.Final",
		"4.1.125.Final",
		"4.1.126.Final",
		"4.1.127.Final",
		"4.1.128.Final",
		"4.1.129.Final",
		"4.1.12.Final",
		"4.1.130.Final",
		"4.1.131.Final",
		"4.1.13.Final",
		"4.1.14.Final",
		"4.1.15.Final",
		"4.1.16.Final",
		"4.1.17.Final",
		"4.1.18.Final",
		"4.1.19.Final",
		"4.1.1.Final",
		"4.1.20.Final",
		"4.1.21.Final",
		"4.1.22.Final",
		"4.1.23.Final",
		"4.1.24.Final",
		"4.1.25.Final",
		"4.1.26.Final",
		"4.1.27.Final",
		"4.1.28.Final",
		"4.1.29.Final",
		"4.1.2.Final",
		"4.1.30.Final",
		"4.1.31.Final",
		"4.1.32.Final",
		"4.1.33.Final",
		"4.1.34.Final",
		"4.1.35.Final",
		"4.1.36.Final",
		"4.1.37.Final",
		"4.1.38.Final",
		"4.1.39.Final",
		"4.1.3.Final",
		"4.1.40.Final",
		"4.1.41.Final",
		"4.1.42.Final",
		"4.1.43.Final",
		"4.1.44.Final",
		"4.1.45.Final",
		"4.1.46.Final",
		"4.1.47.Final",
		"4.1.48.Final",
		"4.1.49.Final",
		"4.1.4.Final",
		"4.1.50.Final",
		"4.1.51.Final",
		"4.1.52.Final",
		"4.1.53.Final",
		"4.1.54.Final",
		"4.1.55.Final",
		"4.1.56.Final",
		"4.1.57.Final",
		"4.1.58.Final",
		"4.1.59.Final",
		"4.1.5.Final",
		"4.1.60.Final",
		"4.1.61.Final",
		"4.1.62.Final",
		"4.1.63.Final",
		"4.1.64.Final",
		"4.1.65.Final",
		"4.1.66.Final",
		"4.1.67.Final",
		"4.1.68.Final",
		"4.1.69.Final",
		"4.1.6.Final",
		"4.1.70.Final",
		"4.1.71.Final",
		"4.1.72.Final",
		"4.1.73.Final",
		"4.1.74.Final",
		"4.1.75.Final",
		"4.1.76.Final",
		"4.1.77.Final",
		"4.1.78.Final",
		"4.1.79.Final",
		"4.1.7.Final",
		"4.1.80.Final",
		"4.1.81.Final",
		"4.1.82.Final",
		"4.1.83.Final",
		"4.1.84.Final",
		"4.1.85.Final",
		"4.1.86.Final",
		"4.1.87.Final",
		"4.1.88.Final",
		"4.1.89.Final",
		"4.1.8.Final",
		"4.1.90.Final",
		"4.1.91.Final",
		"4.1.92.Final",
		"4.1.93.Final",
		"4.1.94.Final",
		"4.1.95.Final",
		"4.1.96.Final",
		"4.1.97.Final",
		"4.1.98.Final",
		"4.1.99.Final",
		"4.1.9.Final",
		"4.2.0.Alpha1",
		"4.2.0.Alpha2",
		"4.2.0.Alpha3",
		"4.2.0.Alpha4",
		"4.2.0.Alpha5",
		"4.2.0.Beta1",
		"4.2.0.Final",
		"4.2.0.RC1",
		"4.2.0.RC2",
		"4.2.0.RC3",
		"4.2.0.RC4",
		"4.2.10.Final",
		"4.2.1.Final",
		"4.2.2.Final",
		"4.2.3.Final",
		"4.2.4.Final",
		"4.2.5.Final",
		"4.2.6.Final",
		"4.2.7.Final",
		"4.2.8.Final",
		"4.2.9.Final",
		"5.0.0.Alpha2",
	}
)
