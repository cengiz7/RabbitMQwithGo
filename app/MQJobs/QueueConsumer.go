package MQJobs

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"

	"../ErrorCheck"
)

func ConsumeFromQueue( queueName string, mqConnection string, indexChannel chan bool, messageCh chan string){
	ch := OpenAmqpChannel(mqConnection)
	err := ch.Qos(
		1,    // prefetch count
		0,     // prefetch size
		false, 	 // global
	)
	ErrorCheck.Check(err, "Failed to set QoS","QoS successfully set for queue: "+queueName)

	queue, err := ch.Consume(
		queueName, 	 	// queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		QueueArgs,    // args
	)
	ErrorCheck.Check(err, "Failed to register a consumer","Consumer successfully set for queue: "+queueName)

	go StartConsume(queue,indexChannel, messageCh)
}


// using <-chan rather that chan amqp... is important
// |<- means| read only channel
func StartConsume(queue <-chan amqp.Delivery, indexChannel <-chan bool, messageCh chan <- string){
	fmt.Println("\n\nQueue consuming started...")
	for d := range queue {
		<- indexChannel
		// notify rabbitmq that the message is consumed
		messageCh <- string(d.Body)
		err := d.Ack(false); if err != nil {
			log.Println("Failed to ack message!!!")
		}
	}
}
