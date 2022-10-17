package utils

const (
	ErrorParse          string = "Failed to parse request body."
	ErrorAcctName       string = "Invalid `name`."
	ErrorAcctEmail      string = "Invalid `email`."
	ErrorAcctPW         string = "Invalid `password`."
	ErrorNonMatchPW     string = "Non-matching password inputs."
	ErrorCreateUser     string = "Failed create user transaction."
	ErrorFailedDB       string = "Failed DB operation."
	ErrorNoRowsAffected string = "result.RowsAffected == 0"
)
