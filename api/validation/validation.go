package validation

import (
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/chirag1807/task-management-system/api/model/response"
	errorhandling "github.com/chirag1807/task-management-system/error"
	validate "github.com/guptaaashutosh/go_validate"
)

func ValidateParameters(r *http.Request, requestBody interface{}, requestParametersMap *map[string]string, requestParametersFiltersMap *map[string]string,
	queryParametersMap *map[string]string, queryParametersFiltersMap *map[string]string, urlParamsError *[]response.InvalidParameters) (error, error, error) {

	var invalidParamErr []response.InvalidParameters
	var invalidParamKeysArr []string

	if urlParamsError != nil {
		invalidParamErr = append(invalidParamErr, *urlParamsError...)
	}

	queryParameterData := validate.FromURLValues(r.URL.Query())
	queryParamData := queryParameterData.Create()

	var queryParameterMapDereference map[string]string
	if queryParametersMap != nil {
		queryParameterMapDereference = *queryParametersMap
		queryParamData.StringRules(queryParameterMapDereference)
	}

	if queryParametersFiltersMap != nil {
		queryParamData.FilterRules(*queryParametersFiltersMap)
	}

	if !queryParamData.Validate() {
		// Range over query parameter map to get particular parameter error
		for key := range queryParameterMapDereference {
			if len(queryParamData.Errors.FieldOne(key)) != 0 {
				invalidParamErr = append(invalidParamErr, response.InvalidParameters{ParameterName: key, ErrorMessage: queryParamData.Errors.FieldOne(key)})
				invalidParamKeysArr = append(invalidParamKeysArr, key)
			}
		}
	}

	// Collection of request parameter
	requestParameterData, err := validate.FromRequest(r)
	if err != nil {
		return err, nil, nil
	}

	requestBodyData := requestParameterData.Create()

	var requestParametersMapDereference map[string]string
	if requestParametersMap != nil {
		requestParametersMapDereference = *requestParametersMap
		requestBodyData.StringRules(requestParametersMapDereference)
	}

	if requestParametersFiltersMap != nil {
		requestBodyData.FilterRules(*requestParametersFiltersMap)
	}

	if !requestBodyData.Validate() {
		// Range over request parameter map to get particular parameter error
		for key := range requestParametersMapDereference {
			if len(requestBodyData.Errors.FieldOne(key)) != 0 {
				invalidParamErr = append(invalidParamErr, response.InvalidParameters{ParameterName: key, ErrorMessage: requestBodyData.Errors.FieldOne(key)})
				invalidParamKeysArr = append(invalidParamKeysArr, key)
			}
		}
	}

	var httpErrorCode int
	if len(invalidParamErr) == 0 {
		if httpErrorCode, err = queryParamData.BindSafeData(requestBody); err != nil {
			queryParamBindDataError := errorhandling.CreateCustomError(err.Error(), httpErrorCode)
			return queryParamBindDataError, nil, nil
		}

		if httpErrorCode, err = requestBodyData.BindSafeData(requestBody); err != nil {
			requestDataBindError := errorhandling.CreateCustomError(err.Error(), httpErrorCode)
			return requestDataBindError, nil, nil
		}
		return nil, nil, nil
	} else {
		invalidParamsSingleLineErrMsg := errorhandling.CreateCustomError(fmt.Sprintf("Invalid data in:%s", strings.Join(invalidParamKeysArr, ", ")), http.StatusBadRequest)

		var error string
		for _, v := range invalidParamErr {
			error += v.ErrorMessage
		}
		invalidParamsMultiLineErrMsg := errorhandling.CreateCustomError(error, http.StatusBadRequest)
		return nil, invalidParamsMultiLineErrMsg, invalidParamsSingleLineErrMsg
	}

}

func EmailValidation(email string) bool {
	emailAddress, err := mail.ParseAddress(email)
	return err == nil && emailAddress.Address == email
}
