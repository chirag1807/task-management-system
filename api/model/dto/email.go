package dto

// Email model info
// @Description Email information with email id of receiver, subject of the email and body of the email.
type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
