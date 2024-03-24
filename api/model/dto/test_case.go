package dto

type ExpectedMessage struct {
	Message    string `json:"message,omitempty"`
	StatusCode *int   `json:"statuscode,omitempty"`
}
