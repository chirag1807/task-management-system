package errorhandling

import (
	"net/http"
)

type CustomError struct {
	StatusCode   int    `json:"statuscode" example:"000"`
	ErrorMessage string `json:"error" example:"Corresponding Error Message will Show Here"`
}

// here I implemented error interface's Error() method. so that we can customize error for our project.
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
	AccessTokenExpired               = CreateCustomError("Access Token is Expired, Please Regenrate It.", http.StatusUnauthorized)
	DuplicateEmailFound              = CreateCustomError("Duplicate Email Found.", http.StatusConflict)
	EmailvalidationError             = CreateCustomError("Email Validation Failed, Please Provide Valid Email.", http.StatusBadRequest)
	LeftAllTeamsToMakeProfilePrivate = CreateCustomError("You must Left All Teams that You are Part of to Make Your Profile Private.", http.StatusUnauthorized)
	MemberExist                      = CreateCustomError("Member Already Added in Team.", http.StatusConflict)
	NoUserFound                      = CreateCustomError("No User Found for This Request.", http.StatusNotFound)
	NoEmailFound                     = CreateCustomError("No User Registered with This Email ID.", http.StatusNotFound)
	NoOTPIDFound                     = CreateCustomError("No OTP ID Found.", http.StatusNotFound)
	NoTaskFound                      = CreateCustomError("No Task Found For This Request.", http.StatusNotFound)
	NotAllowed                       = CreateCustomError("You are not Allowed to Perform this Task.", http.StatusForbidden)
	NotAMember                       = CreateCustomError("You can not Left the Meeting Because You are Not a Member of This Team.", http.StatusUnauthorized)
	OTPVerificationTimeExpired       = CreateCustomError("Sorry, Time for OTP Verification has expired.", http.StatusForbidden)
	OTPNotMatched                    = CreateCustomError("You have Entered Wrong OTP, Try Again with Correct OTP.", http.StatusUnauthorized)
	OnlyPublicMemberAllowed          = CreateCustomError("Only Public Profile Users can be Added in Team.", http.StatusBadRequest)
	OnlyPublicUserAssignne           = CreateCustomError("Tasks can be Assgined to Only Public Profile Users.", http.StatusBadRequest)
	OnlyPublicTeamAssignne           = CreateCustomError("Tasks can be Assgined to Only Public Profile Teams.", http.StatusBadRequest)
	PasswordNotMatched               = CreateCustomError("Password is Incorrect.", http.StatusUnauthorized)
	ProvideValidFlag                 = CreateCustomError("Please Provide Valid Flag to Proceed Further. Flag Value must be either 0 or 1", http.StatusUnprocessableEntity)
	ProvideValidParams               = CreateCustomError("Please Provide Valid URL Parameter to Proceed Further.", http.StatusUnprocessableEntity)
	ReadBodyError                    = CreateCustomError("Could not Read Request Body, Please Provide Valid Body.", http.StatusBadRequest)
	ReadQueryParamsError             = CreateCustomError("Could not Read Request Query Parameters, Please Provide Valid Query Parameters.", http.StatusBadRequest)
	ReadDataError                    = CreateCustomError("Could not Decode the Data, Please Provide Valid Data.", http.StatusBadRequest)
	RegistrationFailedError          = CreateCustomError("User Registration Failed.", http.StatusInternalServerError)
	RefreshTokenExpired              = CreateCustomError("Access Token is Expired, Please Do Login Again.", http.StatusUnauthorized)
	RefreshTokenError                = CreateCustomError("Access Token Can't be Regenerated, Please Do Login Again.", http.StatusUnauthorized)
	RefreshTokenNotFound             = CreateCustomError("Refresh Token Not Found.", http.StatusUnauthorized)
	TokenNotFound                    = CreateCustomError("Authorization Token Not Found.", http.StatusUnauthorized)
	TaskClosed                       = CreateCustomError("Task Can't be Updated because It is Closed.", http.StatusUnprocessableEntity)
)
