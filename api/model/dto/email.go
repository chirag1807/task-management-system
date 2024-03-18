package dto

type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type EmailBody struct {
	Name    string
	Message string
}
