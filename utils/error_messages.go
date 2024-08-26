package utils

const (
	ErrorParse          			string = "Failed to parse request body."
	ErrorAcctName       			string = "Invalid `name`."
	ErrorAcctEmail      			string = "Invalid `email`."
	ErrorAcctPW         			string = "Invalid `password`."
	ErrorAcctPhone						string = "Invalid `phone_number`."
	ErrorNonMatchPW     			string = "Non-matching password inputs."
	ErrorCreateUser     			string = "Failed create user transaction."
	ErrorVaultsCreateUser			string = "Failed vaults API create_user."
	ErrorVaultsCreateVault		string = "Failed vaults API create_vault."
	ErrorVaultsListVaults			string = "Failed vaults API list_vaults."
	ErrorVaultsRetrieveVault	string = "Failed vaults API retrieve_vault."
	ErrorVaultsUpdateVault		string = "Failed vaults API update_vault."
	ErrorVaultsDeleteVault		string = "Failed vaults API delete_vault."
	ErrorVaultsCreateEntry		string = "Failed vaults API create_entry."
	ErrorVaultsRetrieveEntry	string = "Failed vaults API retrieve_entry."
	ErrorVaultsDeleteUser			string = "Failed vaults API delete_user."
	ErrorFailedDB       			string = "Failed DB operation."
	ErrorNoRowsAffected 			string = "result.RowsAffected == 0"
	ErrorIPMismatch 					string = "Different IP addresses."
	ErrorBadClient						string = "Client did something it shouldn't have."
	ErrorParams								string = "Invalid URL parameters."
	ErrorUserContext					string = "Invalid user context."
)
