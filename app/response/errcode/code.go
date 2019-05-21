package errcode

import "sync"

// Error Codes
const (
	ValidationError          = "VALIDATION_ERROR"
	InvalidRequest           = "INVALID_REQUEST"
	SystemError              = "SYSTEM_ERROR"
	APIEndpointNotExist      = "API_ENDPOINT_NOT_EXIST"
	PasswordNotStrength      = "PASSWORD_NOT_STRENGTH"
	PhoneNumberAlreadyExist  = "PHONE_NUMBER_ALREADY_EXIST"
	UserFailedAuthentication = "USER_FAILED_AUTHENTICATION"
	TokenInvalid             = "TOKEN_INVALID"
)

// Message :
var Message sync.Map

func init() {
	Message.Store(ValidationError, "Validation error")
	Message.Store(InvalidRequest, "Request input is invalid")
	Message.Store(SystemError, "System busy, please try again")
	Message.Store(APIEndpointNotExist, "API endpoint not exist")
	Message.Store(PasswordNotStrength, "The password is not strong enough")
	Message.Store(PhoneNumberAlreadyExist, "The phonenumber is already exist")
	Message.Store(UserFailedAuthentication, "The user failed to authenticate")
	Message.Store(TokenInvalid, "Token invalid")
}
