package main

import (
	"time"

	"./HttpJobs"
	"./MQJobs"
)

func main(){
	infinite     := make(chan bool)
	indexChannel := make(chan bool)
	messageCh    := make(chan string)
	connection   := "amqp://guest:guest@localhost:5672/"
	queueName 	 := "guyruk"
	maxPriority  := int64(3)
	// Declare Queue
	MQJobs.DeclarePriorityQueue(queueName, maxPriority, connection)

	// Consume from queue
	MQJobs.ConsumeFromQueue(queueName, connection, indexChannel, messageCh)

	// Set http handlers
	go HttpJobs.SetHandlers(queueName, connection, indexChannel, messageCh)
	time.Sleep( 2 * time.Second)

	// Post sample messages
	HttpJobs.Post()

	<- infinite
}


