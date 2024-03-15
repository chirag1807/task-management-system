package response

type InvalidParameters struct {
	ParameterName string `json:"parameterName"`
	ErrorMessage  string `json:"errorMessage"`
}
