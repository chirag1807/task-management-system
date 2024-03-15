package utils

import (
	"log"
	"net/smtp"

	"github.com/chirag1807/task-management-system/api/model/dto"
	"github.com/chirag1807/task-management-system/config"
)

//SendEmail uses go's built in package net/smtp to send email to given address.
func SendEmail(email dto.Email) {
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
		return
	}
	log.Println("Email Sent Succesfully.")
}
