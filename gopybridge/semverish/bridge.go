// go_shim.go
package main

/*
   #include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"unsafe"

	"github.com/rng70/versions/resolver"
	"github.com/rng70/versions/semver"
	"github.com/rng70/versions/vars"
)

//export NaturalSortedVersions
func NaturalSortedVersions(versions *C.char, descending C._Bool, safeParse C._Bool) *C.char {
	in := C.GoString(versions)
	var _versions []string
	if err := json.Unmarshal([]byte(in), &_versions); err != nil {
		return C.CString("ERROR: invalid json: " + err.Error())
	}
	res := semver.SortedVersions(_versions, bool(descending), bool(safeParse))
	out, _ := json.Marshal(res)
	return C.CString(string(out))
}

//export AnalyzeConstrains
func AnalyzeConstrains(language *C.char, constraints *C.char, versions *C.char) *C.char {
	_language := C.GoString(language)

	in := C.GoString(versions)
	var _versions []string
	if err := json.Unmarshal([]byte(in), &_versions); err != nil {
		return C.CString("ERROR: invalid json: " + err.Error())
	}

	_constraints := C.GoString(constraints)

	style := vars.StylePy
	switch _language {
	case "python", "py":
		style = vars.StylePy
	case "nuget", "csharp", "dotnet", "cs":
		style = vars.StyleNuGet
	case "npm", "node", "nodejs", "javascript", "js":
		style = vars.StyleNPM
	case "maven", "java":
		style = vars.StyleNPM
	default:
		style = vars.StylePy
	}

	analysis := resolver.AnalyzeConstraint(style, _constraints, _versions)

	out, _ := json.Marshal(analysis.Matches)

	return C.CString(string(out))
}

//export FreeCString
func FreeCString(str *C.char) {
	C.free(unsafe.Pointer(str))
}

func main() {}
