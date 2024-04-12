package utils

import (
	"encoding/json"
	"fmt"

	"github.com/chirag1807/task-management-system/api/model/dto"
	"github.com/chirag1807/task-management-system/constant"
	amqp "github.com/rabbitmq/amqp091-go"
)

// ConsumeUserMail function make use of rabbitmq and consume messgae from default queue.
func ConsumeUserMail(rabbitmqConn *amqp.Connection) {
	ch, err := rabbitmqConn.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		constant.USER_MAIL_QUEUE,
		constant.EMPTY_STRING,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	var forever = make(chan bool)

	go func() {
		for m := range msgs {
			var email dto.Email
			_ = json.Unmarshal(m.Body, &email)
			fmt.Println(email.To, email.Subject, email.Body)
			SendEmail(email)
		}
	}()

	<-forever
	// <-forever which blocks our main function from completing until the channel is satisfied.
}
