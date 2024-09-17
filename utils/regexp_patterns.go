package utils

import "regexp"

var (
	EmailRegexp			 			 = regexp.MustCompile("^[\u00BF-\u1FFF\u2C00-\uD7FF\\w\\.\\+\\-]+@[\u00BF-\u1FFF\u2C00-\uD7FF\\w\\-]+(\\.[\u00BF-\u1FFF\u2C00-\uD7FF\\w\\-]+)+$")
	NameRegexp			 			 = regexp.MustCompile(`^[\\u00BF-\\u1FFF\\u2C00-\\uD7FF\w \-\,\.\']{1,50}$`)
	PhoneRegexp			 			 = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	HexEncodedKeyRegexp		 = regexp.MustCompile(`^[A-Fa-f0-9]{64}$`)
	SlugRegexp       			 = regexp.MustCompile(`^[\w-]{32}$`)
	TokenRegexp      			 = regexp.MustCompile(`^[\w-]{80}$`)
	RowsRegexp       			 = regexp.MustCompile(`^result.RowsAffected \([0-9]+\) > 1$`)
	AuthHeaderRegexp 			 = regexp.MustCompile(`^[Tt]oken [\w-]{80}$`)
	TokenNullRegexp  			 = regexp.MustCompile(`^[Tt]oken (null)?$`)
	UniqueConstraintRegexp = regexp.MustCompile(`^(UNIQUE constraint failed: users\.(email_address|phone_number)|ERROR: duplicate key value violates unique constraint "users_(email_address|phone_number)_key" \(SQLSTATE 23505\))$`)
)
