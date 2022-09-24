package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// https://www.youtube.com/watch?v=7rkeORD4jSw
// https://en.wikipedia.org/wiki/Advanced_Message_Queuing_Protocol
// decouple
// scalable
// performent
// flexibility
// cross programming languages

// term
// message: data that need to get transmitted (data, file, metadata,...)
// consumer: the service that using messages from queue
/* AMQP parts
   Exchange: Receive and distribute the messages
   Queue: store messages from exchange
   binding: Connection between the service sending messages and the exchange
   binding key: the string refernce to the exchange
*/
// acks: the messages in the queue remove only when the consumer let the broker -> prevent loss messages
// subcribe to the queue: listening for messages
// fanout: 1 sending to all consumers
// direct exchange: compare the key in the message and the binding key, if match --> sending to that consumers
// topic exchange: send to consumer that match with the topic (compare to the direct exchange, this one is sending to more consumers)
// header exchange: ignore the message, check the header for sending
// default exchange (nameless exchange): send to the queue with the specific name
func main() {
	// connect RabbitMQ
	log.Println("Connecting to rabbitmq")
	rabbitmqConn, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		close(rabbitmqConn)
	}()
	// start listening
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create consummers
	consumer, err := event.NewConsummer(rabbitmqConn)
	if err != nil {
		log.Fatal(err)
	}

	// watch queues and consume events (topics)
	err = consumer.Listen([]string{infoLog, warningLog, errLog})
	if err != nil {
		log.Fatal(err)
	}
}

const (
	infoLog    = "log.INFO"
	warningLog = "log.WARNING"
	errLog     = "log.ERROR"
)

func close(rabbitmqConn *amqp.Connection) {
	err := rabbitmqConn.Close()
	if err != nil {
		// Server may not receive close signal
		var retryingTime = 5
		for err != nil {
			time.Sleep(time.Second * 5)
			err = nil
			retryingTime--
			if retryingTime < 0 {
				return
			}
			err = rabbitmqConn.Close()
		}
	}
}

func connect() (*amqp.Connection, error) {
	var count int64
	var backOff = time.Second
	var con *amqp.Connection

	// do not continue until rabbitmq ready
	for con == nil {
		count++
        // "rabbitmq" service in docker compose
		c, err := amqp.Dial(fmt.Sprintf("amqp://%v:%v@%v", "guest", "guest", "rabbitmq"))
		if err != nil {
			log.Printf("RabbitMQ not yet ready, %v time with error: %v", count, err)
		} else {
			con = c
		}

		if count > 5 {
			return nil, fmt.Errorf("exceed %v time for connecting RabbitMQ", count)
		}

		backOff = time.Second * time.Duration(math.Pow(float64(count), 2))
		log.Printf("Sleep %v", backOff.Seconds())
		time.Sleep(backOff)

	}

	return con, nil
}
