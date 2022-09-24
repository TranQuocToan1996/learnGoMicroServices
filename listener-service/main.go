package main

import (
	"fmt"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// connect RabbitMQ
	rabbitmqConn, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		close(rabbitmqConn)
	}()
	log.Println("Connecting to rabbitmq")

	// start listening

	// create consummers

	// watch queues and consume events (topics)
}

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
		c, err := amqp.Dial(fmt.Sprintf("amqp://%v:%v@%v", "guest", "guest", "localhost"))
		if err != nil {
			log.Printf("RabbitMQ not yet ready, %v time with error: %v", count, err)
		} else {
			con = c
		}

		if count > 5 {
			return nil, fmt.Errorf("exceed %v time for connecting RabbitMQ", count)
		}

		backOff = time.Second * time.Duration(math.Pow(float64(count, 2)))
		log.Printf("Sleep %v", backOff.Seconds())
		time.Sleep(backOff)

	}

	return con, nil
}
