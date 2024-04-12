package errorhandling

import (
	"net/http"

	"github.com/chirag1807/task-management-system/constant"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type InvalidRequestData struct {
	ParameterName string `json:"param" example:"email"`
	ErrorMessage  string `json:"error" example:"please provide email in valid format."`
}

type RequestDataValidationError struct {
	StatusCode int                  `json:"statuscode" example:"000"`
	Errors     []InvalidRequestData `json:"errors"`
}

// here I implemented error interface's Error() method. so that we can customize request data validation errors for our project.
func (r RequestDataValidationError) Error() string {
	var error string
	for _, v := range r.Errors {
		error += v.ErrorMessage + ", "
	}
	return error
}

func CreateRequestDataValidationError(errors []InvalidRequestData, statusCode int) error {
	return RequestDataValidationError{
		StatusCode: statusCode,
		Errors:     errors,
	}
}

func HandleInvalidRequestData(w http.ResponseWriter, r *http.Request, err error, translator ut.Translator) {
	validationErrors := err.(validator.ValidationErrors)
	var errors []InvalidRequestData

	for _, e := range validationErrors {
		errors = append(errors, InvalidRequestData{ParameterName: e.Field(), ErrorMessage: e.Translate(translator)})
	}

	SendErrorResponse(r, w, CreateRequestDataValidationError(errors, http.StatusBadRequest), constant.EMPTY_STRING)
}
