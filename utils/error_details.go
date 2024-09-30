package utils

const (
	ErrorBadRequest  	string = "Bad request."
	ErrorToken       	string = "Invalid token."
	ErrorMFAToken			string = "Invalid MFA token."
	ErrorPhoneOTP			string = "Invalid phone OTP."
	ErrorServer      	string = "Oops, something went wrong!"
	ErrorDiffEmail   	string = "Oops, failed to create account - try using a different email address or phone number."
	ErrorFailedLogin 	string = "Oops, failed to log in - try again!"
	ErrorAuthenticate	string = "Oops, failed to authenticate - try again!"
)
