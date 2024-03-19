package utils

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"
	"strconv"

	"github.com/chirag1807/task-management-system/api/model/dto"
	"github.com/chirag1807/task-management-system/config"
)

// SendEmail uses go's built in package net/smtp to send email to given address.
func SendEmail(email dto.Email) error {
	to := []string{
		email.To,
	}

	msg := []byte("Subject: " + email.Subject + "\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		email.Body)

	auth := smtp.PlainAuth(
		"",
		config.Config.SMTP.EmailFrom,
		config.Config.SMTP.EmailPassword,
		config.Config.SMTP.Host,
	)

	err := smtp.SendMail(
		config.Config.SMTP.Host+":"+config.Config.SMTP.Port,
		auth,
		config.Config.SMTP.EmailFrom,
		to,
		msg,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Email Sent Succesfully.")
	return nil
}

func ParseTemplate(fileName string, data interface{}) ([]byte, error) {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return []byte{}, err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

func PrepareEmailBody(OTP int) string {
	body := `
    <!DOCTYPE html>
    <html lang="en">

    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Email Verification</title>
    </head>

    <body style="font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f4f4f4;">
        <div style="background-color: #2196F3; color: white; text-align: center; padding: 20px;">
            <h2>ZURU TECH</h2>
        </div>

        <div style="padding: 20px;">
            <p>Hello User,</p>
            <p>Please verify your email address by entering below code.</p>
            <p>Your verification code is: <strong>` + strconv.Itoa(OTP) + `</strong></p>
            <p>If you did not initiate this change, please contact our support team immediately.</p>
            <p>Best regards,<br>ZURU TECH</p>
        </div>
    </body>

    </html>
`

	return body
}
