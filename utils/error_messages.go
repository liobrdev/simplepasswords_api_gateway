package utils

const (
	ErrorParse          	string = "Failed to parse request body."
	ErrorAcctName       	string = "Invalid `name`."
	ErrorAcctEmail      	string = "Invalid `email`."
	ErrorAcctPW         	string = "Invalid `password`."
	ErrorNonMatchPW     	string = "Non-matching password inputs."
	ErrorCreateUser     	string = "Failed create user transaction."
	ErrorVaultsCreateUser string = "Failed vaults API create_user."
	ErrorFailedDB       	string = "Failed DB operation."
	ErrorNoRowsAffected 	string = "result.RowsAffected == 0"
	ErrorIPMismatch 			string = "Different IP addresses."
	ErrorBadClient				string = "Client did something it shouldn't have."
)
