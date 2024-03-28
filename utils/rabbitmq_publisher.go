package utils

import (
	"context"
	"encoding/json"

	"github.com/chirag1807/task-management-system/api/model/dto"
	"github.com/chirag1807/task-management-system/constant"
	amqp "github.com/rabbitmq/amqp091-go"
)

// ProduceEmail function make use of rabbitmq and produce messgae to default queue.
func ProduceEmail(rabbitmqConn *amqp.Connection, userEmail dto.Email) (error) {
	ch, err := rabbitmqConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		constant.USER_MAIL_QUEUE,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	messageToProduce, _ := json.Marshal(userEmail)

	err = ch.PublishWithContext(context.Background(), "", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        messageToProduce,
	})

	if err != nil {
		return err
	}

	go ConsumeUserMail(rabbitmqConn)
	return nil
}
