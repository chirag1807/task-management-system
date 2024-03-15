package errorhandling

import (
	"encoding/json"
	"log"
	"net/http"
)

type CustomError struct {
	StatusCode   int
	ErrorMessage string
}

// here I implemented error interface's Error() method.
// so that we can customize error for our project.
func (c CustomError) Error() string {
	return c.ErrorMessage
}

// CreateCustomError takes error message and http status code as parameters and return error in CustomError format.
func CreateCustomError(errorMessage string, statusCode int) error {
	return CustomError{
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
	}
}

var (
	ReadBodyError           = CreateCustomError("Could not Read Request Body, Please Provide Valid Body.", http.StatusBadRequest)
	ReadDataError           = CreateCustomError("Could not Decode the Data, Please Provide Valid Data.", http.StatusBadRequest)
	EmailvalidationError    = CreateCustomError("Email Validation Failed, Please Provide Valid Email.", http.StatusBadRequest)
	DuplicateEmailFound     = CreateCustomError("Duplicate Email Found.", http.StatusConflict)
	RegistrationFailedError = CreateCustomError("User Registration Failed.", http.StatusInternalServerError)
	LoginFailedError        = CreateCustomError("User Login Failed.", http.StatusUnauthorized)
	AccessTokenExpired      = CreateCustomError("Access Token is Expired, Please Regenrate It.", http.StatusUnauthorized)
	RefreshTokenExpired     = CreateCustomError("Access Token is Expired, Please Do Login Again.", http.StatusUnauthorized)
	RefreshTokenError       = CreateCustomError("Access Token Can't be Regenerated, Please Do Login Again.", http.StatusUnauthorized)
	UnauthorizedError       = CreateCustomError("You are Not Authorized to Perform this Action.", http.StatusUnauthorized)
	NoUserFound             = CreateCustomError("No User Found for This Request.", http.StatusNotFound)
	RefreshTokenNotFound    = CreateCustomError("Refresh Token Not Found.", http.StatusUnauthorized)
	PasswordNotMatch        = CreateCustomError("Password is Incorrect.", http.StatusUnauthorized)
)

// SendErrorResponse send defined errors in response with error message and status code.
// and for those errors, which are not defined in global error handling,
// it will simply send 'Internal Server Error' as error message and 500 as status code.
func SendErrorResponse(w http.ResponseWriter, err error) {
	var response interface{}

	log.Println(err.Error())

	if error, ok := err.(CustomError); ok {
		response = CustomError{
			StatusCode: error.StatusCode,
			ErrorMessage:    error.ErrorMessage,
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(error.StatusCode)

	} else {
		response = CustomError{
			StatusCode: http.StatusInternalServerError,
			ErrorMessage:    "Internal Server Error",
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(response)
}
