package MQJobs

import (
	"github.com/streadway/amqp"

	"investment/app/ErrorCheck"
)

func SentToQueue(message []byte, queueName,connection string, priority uint8) {
	ch := OpenAmqpChannel(connection)
	err := ch.Publish(
		"",     // exchange
		queueName, 		// routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            message ,
			DeliveryMode:    amqp.Persistent,
			Priority:        priority,
		})
	ErrorCheck.Check(err,"Error while sending message to "+queueName, "Message succesfully sent to "+queueName)
}