package response

// InvalidParameters model info
// @Description Invalid parameter with name and corresponding error message.
type InvalidParameters struct {
	ParameterName string `json:"parameterName" example:"email"`
	ErrorMessage  string `json:"errorMessage" example:"please provide email in valid format."`
}
