package utils

const (
	AuthFirstFactor		 string = "auth_first_factor"
	AuthSecondFactor	 string = "auth_second_factor"
	CreateAccount			 string = "create_account"
	LogoutAccount			 string = "logout_account"
	TestAuthReq				 string = "test_auth_req"
	RetrieveUser			 string = "retrieve_user"
	VerifyEmailTry		 string = "verify_email_try"
	VerifyEmailConfirm string = "verify_email_confirm"
	VerifyPhoneTry		 string = "verify_phone_try"
	VerifyPhoneConfirm string = "verify_phone_confirm"

	// vaults
	CreateUser    string = "create_user"
	CreateVault   string = "create_vault"
	CreateEntry   string = "create_entry"
	CreateSecret  string = "create_secret"
	ListVaults		string = "list_vaults"
	RetrieveVault string = "retrieve_vault"
	RetrieveEntry string = "retrieve_entry"
	UpdateVault   string = "update_vault"
	UpdateEntry   string = "update_entry"
	UpdateSecret  string = "update_secret"
	MoveSecret		string = "move_secret"
	DeleteVault   string = "delete_vault"
	DeleteEntry   string = "delete_entry"
	DeleteSecret  string = "delete_secret"
)
