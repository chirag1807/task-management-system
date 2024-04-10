package errorhandling

import (
	"encoding/json"
	"net/http"

	"github.com/chirag1807/task-management-system/config"
)

// SendErrorResponse send defined errors in response with error message and status code.
// and for those errors, which are not defined in global error handling,
// it will simply send 'Internal Server Error' as error message and 500 as status code.
func SendErrorResponse(r *http.Request, w http.ResponseWriter, err error, message string, params ...interface{}) {
	if error, ok := err.(CustomError); ok {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(error.StatusCode)
		config.LoggerInstance.Warning(err.Error())

	} else if error, ok := err.(RequestDataValidationError); ok {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(error.StatusCode)
		config.LoggerInstance.Warning(err.Error())
	} else {
		config.LoggerInstance.Error(r, err, message, params...)
		err = CustomError{
			StatusCode:   http.StatusInternalServerError,
			ErrorMessage: "Internal Server Error",
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(err)
}
