package main

import (
	"brokerservice/model"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	log.Println("Connecting to rabbitmq")
	rabbitmqConn, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		close(rabbitmqConn)
	}()

	// cfg := model.GetConfig()
	cfg := &model.Config{
		Port:     "81",
		RabbitMQ: rabbitmqConn,
	}

	mux := routes(cfg)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: mux,
	}

	log.Printf("Starting on port %v", cfg.Port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

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
