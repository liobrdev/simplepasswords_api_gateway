package utils

import "regexp"

var (
	EmailRegexp			 = regexp.MustCompile("^[\u00BF-\u1FFF\u2C00-\uD7FF\\w\\.\\+\\-]+@[\u00BF-\u1FFF\u2C00-\uD7FF\\w\\-]+(\\.[\u00BF-\u1FFF\u2C00-\uD7FF\\w\\-]+)+$")
	NameRegexp			 = regexp.MustCompile("^[\u00BF-\u1FFF\u2C00-\uD7FF\\w\\s\\-\\,\\.\\']{1,50}$")
	SlugRegexp       = regexp.MustCompile(`^[\w-]{32}$`)
	TokenRegexp      = regexp.MustCompile(`^[\w-]{80}$`)
	RowsRegexp       = regexp.MustCompile(`^result.RowsAffected \([0-9]+\) > 1$`)
	AuthHeaderRegexp = regexp.MustCompile(`^[Tt]oken [\w-]{80}$`)
	TokenNullRegexp  = regexp.MustCompile(`^[Tt]oken (null)?$`)
)
