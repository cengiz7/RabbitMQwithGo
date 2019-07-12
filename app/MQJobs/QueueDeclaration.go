package MQJobs

import (
	"github.com/streadway/amqp"
	"log"
	"strconv"

	"../ErrorCheck"
)

// Shared queue args variable
var QueueArgs = make(amqp.Table)

func ConnectAmqp(mqConnection string) *amqp.Connection {
	conn, err := amqp.Dial(mqConnection)
	ErrorCheck.Check(err, "Failed to connect RabbitMQ","Succesfully connected to RabbitMQ")
	return conn
}

func OpenAmqpChannel(mqConnection string) *amqp.Channel {
	conn := ConnectAmqp(mqConnection)
	ch, err := conn.Channel()
	ErrorCheck.Check(err, "Failed to open a channel","Successfully opened a amqp channel.")
	return ch
}

func DeclarePriorityQueue(queueName string, maxPriority int64, mqConnection string ){
	QueueArgs["x-max-priority"] = maxPriority
	ch := OpenAmqpChannel(mqConnection)
	resp, err := ch.QueueDeclare(
		queueName, // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		QueueArgs,     // arguments
	)
	ErrorCheck.Check(err, "Failed to declare queue","Queue successfully declared on RabbitMQ.")
	log.Println("\nQueue currently has "+strconv.Itoa(resp.Messages)+" messages.\n")
	ch.Close()
}