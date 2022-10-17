package utils

import "regexp"

var (
	SlugRegexp  = regexp.MustCompile(`^[\w-]{32}$`)
	TokenRegexp = regexp.MustCompile(`^[\w-]{144}$`)
	RowsRegexp  = regexp.MustCompile(`^result.RowsAffected \([0-9]+\) > 1$`)
)
