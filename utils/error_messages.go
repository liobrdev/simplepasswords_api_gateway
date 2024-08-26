package utils

const (
	ErrorParse          		string = "Failed to parse request body."
	ErrorAcctName       		string = "Invalid `name`."
	ErrorAcctEmail      		string = "Invalid `email`."
	ErrorAcctPW         		string = "Invalid `password`."
	ErrorAcctPhone					string = "Invalid `phone_number`."
	ErrorNonMatchPW     		string = "Non-matching password inputs."
	ErrorCreateUser     		string = "Failed create user transaction."
	ErrorVaultsCreateUser		string = "Failed vaults API create_user."
	ErrorVaultsDeleteUser		string = "Failed vaults API delete_user."
	ErrorVaultsListVaults		string = "Failed vaults API list_vaults."
	ErrorFailedDB       		string = "Failed DB operation."
	ErrorNoRowsAffected 		string = "result.RowsAffected == 0"
	ErrorIPMismatch 				string = "Different IP addresses."
	ErrorBadClient					string = "Client did something it shouldn't have."
	ErrorParams							string = "Invalid URL parameters."
	ErrorUserContext				string = "Invalid user context."
)
