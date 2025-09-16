/********************************************************************************************************
*                                                                                                       *
* OPENREFACTORY CONFIDENTIAL                                                                            *
* --------------------------                                                                            *
*                                                                                                       *
* Copyright (c) 2025 OpenRefactory, Inc. All Rights Reserved.                                           *
*                                                                                                       *
* NOTICE: All information contained herein is, and remains the property of OpenRefactory, Inc. The      *
* intellectual and technical concepts contained herein are proprietary to OpenRefactory, Inc. and       *
* may be covered by U.S. and Foreign Patents, patents in process, and are protected by trade secret     *
* or copyright law. Dissemination of this information or reproduction of this material is strictly      *
* forbidden unless prior written permission is obtained from OpenRefactory, Inc.                        *
*                                                                                                       *
* Contributors: Al Arafat Tanin (OpenRefactory, Inc.)                                                   *
*********************************************************************************************************/

package resolver

import "regexp"

type Constraint struct {
	Op  string // =, !=, <, <=, >, >=, ^, ~
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
)

type ConstraintResult struct {
	Raw     string
	Parsed  string
	Matches []string
}

// Notes: Go's regexp does not support non-capturing groups (?:...), so parentheses are used.
// We capture two main groups: operator (group 1) and token (group 2).
var (
	reDashRange  = regexp.MustCompile(`^\s*([0-9]+(?:\.[0-9]+){2})\s*-\s*([0-9]+(?:\.[0-9]+){2})\s*$`)
	reNpmToken   = regexp.MustCompile(`(?i)(<=|>=|<|>|=|~|\^)?\s*([0-9]+(\.[0-9]+){0,2}(\.(x|X|\*))?|latest|npm:[^\s@]+@\d+(\.\d+){0,2}|https?://\S+|file:[^\s]+|\*|[0-9]+)`)
	rePyPart     = regexp.MustCompile(`^(==|!=|<=|>=|<|>|~=|===)\s*([0-9]+(\.[0-9]+){0,2}(\.\*)?)\s*$`)
	reNuGetRange = regexp.MustCompile(`^\s*([\[\(])\s*([^,\s]*)\s*,\s*([^\]\)\s]*)\s*([\]\)])\s*$`)
)
