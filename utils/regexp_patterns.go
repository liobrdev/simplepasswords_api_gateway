package utils

import "regexp"

var (
	SlugRegexp       = regexp.MustCompile(`^[\w-]{32}$`)
	TokenRegexp      = regexp.MustCompile(`^[\w-]{80}$`)
	RowsRegexp       = regexp.MustCompile(`^result.RowsAffected \([0-9]+\) > 1$`)
	AuthHeaderRegexp = regexp.MustCompile(`^[Tt]oken [\w-]{80}$`)
	TokenNullRegexp  = regexp.MustCompile(`^[Tt]oken (null)?$`)
)
