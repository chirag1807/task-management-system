package errorhandling

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chirag1807/task-management-system/constant"
	"github.com/gorilla/schema"
)

type CustomError struct {
	ErrorCode      string `json:"code" example:"Bad request"`
	HttpStatusCode int    `json:"-" example:"400"` // json:"-" refers field won't include in response.
	ErrorMessage   string `json:"error" example:"Corresponding Error Message will Show Here"`
}

// here I implemented error interface's Error() method. so that we can customize error for our project.
func (c CustomError) Error() string {
	return c.ErrorMessage
}

// CreateCustomError takes error message and http status code as parameters and return error in CustomError format.
func CreateCustomError(errorMessage string, errorCode string, httpStatusCode int) error {
	return CustomError{
		ErrorCode:      errorCode,
		HttpStatusCode: httpStatusCode,
		ErrorMessage:   errorMessage,
	}
}

var (
	AccessTokenExpired                = CreateCustomError("Access Token is Expired, Please Regenrate It.", http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	DuplicateEmailFound               = CreateCustomError("Duplicate Email Found.", http.StatusText(http.StatusConflict), http.StatusConflict)
	LeftAllTeamsToMakePrivacyPrivate  = CreateCustomError("You must Left All Teams that You are Part of to Make Your Privacy Private.", http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	MemberExist                       = CreateCustomError("Member Already Added in Team.", http.StatusText(http.StatusConflict), http.StatusConflict)
	NoUserFound                       = CreateCustomError("No User Found for This Request.", http.StatusText(http.StatusNotFound), http.StatusNotFound)
	NoEmailFound                      = CreateCustomError("No User Registered with This Email ID.", http.StatusText(http.StatusNotFound), http.StatusNotFound)
	NoOTPIDFound                      = CreateCustomError("No OTP ID Found.", http.StatusText(http.StatusNotFound), http.StatusNotFound)
	NoTaskFound                       = CreateCustomError("No Task Found For This Request.", http.StatusText(http.StatusNotFound), http.StatusNotFound)
	NotAllowed                        = CreateCustomError("You are not Allowed to Perform this Task.", http.StatusText(http.StatusForbidden), http.StatusForbidden)
	NotAMember                        = CreateCustomError("You can not Left the Meeting Because You are Not a Member of This Team.", http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	FirstVerifyOTP                    = CreateCustomError("First Verify OTP with Our System", http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	OTPVerificationTimeExpired        = CreateCustomError("Sorry, Time for OTP Verification has expired.", http.StatusText(http.StatusGone), http.StatusGone)
	OTPNotMatched                     = CreateCustomError("You have Entered Wrong OTP, Try Again with Correct OTP.", http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	OnlyPublicMemberAllowed           = CreateCustomError("Only Public Profile Users can be Added in Team.", http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	OnlyPublicUserAssignne            = CreateCustomError("Tasks can be Assgined to Only Public Profile Users.", http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	OnlyPublicTeamAssignne            = CreateCustomError("Tasks can be Assgined to Only Public Profile Teams.", http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	PasswordNotMatched                = CreateCustomError("Password is Incorrect.", http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	PasswordConfirmPasswordNotMatched = CreateCustomError("Password and Confirm Password Not Matched.", http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	ProvideValidParams                = CreateCustomError("Please Provide Valid URL Parameter to Proceed Further.", http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	ReadBodyError                     = CreateCustomError("Could not Read Request Body, Please Provide Valid Body.", http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	RefreshTokenExpired               = CreateCustomError("Access Token is Expired, Please Do Login Again.", http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	RefreshTokenError                 = CreateCustomError("Access Token Can't be Regenerated, Please Do Login Again.", http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	RefreshTokenNotFound              = CreateCustomError("Refresh Token Not Found.", http.StatusText(http.StatusNotFound), http.StatusNotFound)
	TokenNotFound                     = CreateCustomError("Authorization Token Not Found.", http.StatusText(http.StatusNotFound), http.StatusNotFound)
	TaskClosed                        = CreateCustomError("Task Can't be Updated because It is Closed.", http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
)

// HandleJSONUnmarshalError function handles JSON unmarshalling errors, constructs custom error messages,
// and sends appropriate HTTP error responses.
func HandleJSONUnmarshlError(r *http.Request, w http.ResponseWriter, err error) {
	if syntaxError, ok := err.(*json.SyntaxError); ok {
		syntaxErrorMessage := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		customSyntaxError := CreateCustomError(syntaxErrorMessage, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		SendErrorResponse(r, w, customSyntaxError, constant.EMPTY_STRING)
	} else if unmarshalTypeError, ok := err.(*json.UnmarshalTypeError); ok {
		unmarshalTypeErrorMessage := fmt.Sprintf("Request body contains an invalid value for the field %s", unmarshalTypeError.Field)
		customUnmarshalTypeError := CreateCustomError(unmarshalTypeErrorMessage, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		SendErrorResponse(r, w, customUnmarshalTypeError, constant.EMPTY_STRING)
	} else {
		customError := CreateCustomError(err.Error(), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		SendErrorResponse(r, w, customError, constant.EMPTY_STRING)
	}
}

func HandleSchemaDecodeError(r *http.Request, w http.ResponseWriter, err error) {
	if multiErr, ok := err.(schema.MultiError); ok {
		var fields []string
		for fieldName := range multiErr {
			fields = append(fields, fieldName)
		}
		multiErrorMessage := fmt.Sprintf("Error decoding query parameters for fields: %v", fields)
		customSyntaxError := CreateCustomError(multiErrorMessage, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		SendErrorResponse(r, w, customSyntaxError, constant.EMPTY_STRING)
		return
	} else {
		customError := CreateCustomError(err.Error(), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		SendErrorResponse(r, w, customError, constant.EMPTY_STRING)
	}
}
