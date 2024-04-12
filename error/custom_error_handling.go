package errorhandling

import (
	"net/http"
)

type CustomError struct {
	Code         string `json:"code" example:"Bad request"`
	ErrorMessage string `json:"error" example:"Corresponding Error Message will Show Here"`
}

// here I implemented error interface's Error() method. so that we can customize error for our project.
func (c CustomError) Error() string {
	return c.ErrorMessage
}

// CreateCustomError takes error message and http status code as parameters and return error in CustomError format.
func CreateCustomError(errorMessage string, Code string) error {
	return CustomError{
		Code:         Code,
		ErrorMessage: errorMessage,
	}
}

var (
	AccessTokenExpired                = CreateCustomError("Access Token is Expired, Please Regenrate It.", http.StatusText(http.StatusUnauthorized))
	DuplicateEmailFound               = CreateCustomError("Duplicate Email Found.", http.StatusText(http.StatusConflict))
	EmailvalidationError              = CreateCustomError("Email Validation Failed, Please Provide Valid Email.", http.StatusText(http.StatusBadRequest))
	LeftAllTeamsToMakePrivacyPrivate  = CreateCustomError("You must Left All Teams that You are Part of to Make Your Privacy Private.", http.StatusText(http.StatusUnauthorized))
	MemberExist                       = CreateCustomError("Member Already Added in Team.", http.StatusText(http.StatusConflict))
	NoUserFound                       = CreateCustomError("No User Found for This Request.", http.StatusText(http.StatusNotFound))
	NoEmailFound                      = CreateCustomError("No User Registered with This Email ID.", http.StatusText(http.StatusNotFound))
	NoOTPIDFound                      = CreateCustomError("No OTP ID Found.", http.StatusText(http.StatusNotFound))
	NoTaskFound                       = CreateCustomError("No Task Found For This Request.", http.StatusText(http.StatusNotFound))
	NotAllowed                        = CreateCustomError("You are not Allowed to Perform this Task.", http.StatusText(http.StatusForbidden))
	NotAMember                        = CreateCustomError("You can not Left the Meeting Because You are Not a Member of This Team.", http.StatusText(http.StatusNotFound))
	OTPVerificationTimeExpired        = CreateCustomError("Sorry, Time for OTP Verification has expired.", http.StatusText(http.StatusForbidden))
	OTPNotMatched                     = CreateCustomError("You have Entered Wrong OTP, Try Again with Correct OTP.", http.StatusText(http.StatusUnauthorized))
	OnlyPublicMemberAllowed           = CreateCustomError("Only Public Profile Users can be Added in Team.", http.StatusText(http.StatusBadRequest))
	OnlyPublicUserAssignne            = CreateCustomError("Tasks can be Assgined to Only Public Profile Users.", http.StatusText(http.StatusBadRequest))
	OnlyPublicTeamAssignne            = CreateCustomError("Tasks can be Assgined to Only Public Profile Teams.", http.StatusText(http.StatusBadRequest))
	PasswordNotMatched                = CreateCustomError("Password is Incorrect.", http.StatusText(http.StatusUnauthorized))
	PasswordConfirmPasswordNotMatched = CreateCustomError("Password and Confirm Password Not Matched.", http.StatusText(http.StatusBadRequest))
	ProvideValidFlag                  = CreateCustomError("Please Provide Valid Flag to Proceed Further. Flag Value must be either 0 or 1", http.StatusText(http.StatusUnprocessableEntity))
	ProvideValidParams                = CreateCustomError("Please Provide Valid URL Parameter to Proceed Further.", http.StatusText(http.StatusUnprocessableEntity))
	ReadBodyError                     = CreateCustomError("Could not Read Request Body, Please Provide Valid Body.", http.StatusText(http.StatusBadRequest))
	ReadQueryParamsError              = CreateCustomError("Could not Read Request Query Parameters, Please Provide Valid Query Parameters.", http.StatusText(http.StatusBadRequest))
	ReadDataError                     = CreateCustomError("Could not Decode the Data, Please Provide Valid Data.", http.StatusText(http.StatusBadRequest))
	RegistrationFailedError           = CreateCustomError("User Registration Failed.", http.StatusText(http.StatusInternalServerError))
	RefreshTokenExpired               = CreateCustomError("Access Token is Expired, Please Do Login Again.", http.StatusText(http.StatusUnauthorized))
	RefreshTokenError                 = CreateCustomError("Access Token Can't be Regenerated, Please Do Login Again.", http.StatusText(http.StatusUnauthorized))
	RefreshTokenNotFound              = CreateCustomError("Refresh Token Not Found.", http.StatusText(http.StatusUnauthorized))
	TokenNotFound                     = CreateCustomError("Authorization Token Not Found.", http.StatusText(http.StatusUnauthorized))
	TaskClosed                        = CreateCustomError("Task Can't be Updated because It is Closed.", http.StatusText(http.StatusUnprocessableEntity))
)
