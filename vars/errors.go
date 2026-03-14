package vars

import "errors"

// ErrUnsupportedSource is returned when a constraint string is a URL or file
// path rather than a parseable version constraint.
var ErrUnsupportedSource = errors.New("unsupported version source (URL or file)")
