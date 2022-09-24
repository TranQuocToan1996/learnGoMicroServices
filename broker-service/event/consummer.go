package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	Exchange = "logs_topic"
)

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type Consummer struct {
	conn  *amqp.Connection
	queue string
}

func (c *Consummer) setup() error {
	channel, err := c.conn.Channel()
	// Errors returned from this method will close the channel.
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

func (c *Consummer) Listen(topics []string) error {
	channel, err := c.conn.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	queue, err := declareRandomQueue(channel)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		// bind channel to each topic
		err := channel.QueueBind(
			queue.Name, // random generates when declare queue
			topic,
			Exchange,
			false, // nowait
			nil,
		)

		if err != nil {
			return err
		}
	}

	autoAck, exclusive, noLocal, noWait := true, true, true, true

	messages, err := channel.Consume(queue.Name, "", autoAck, !exclusive, !noLocal, !noWait, nil)
	if err != nil {
		return err
	}

	// Listen until exit this application
	forever := make(chan struct{})
	go func() {
		// Delivery captures the fields for a previously delivered message resident in
		// a queue to be delivered by the server to a consumer from Channel.Consume or
		// Channel.Get.
		for delivery := range messages {
			payload := &Payload{}
			err := json.Unmarshal(delivery.Body, payload)
			if err != nil {
				buf, _ := json.Marshal(delivery)
				fmt.Println(err, string(buf))
			}

			go handlePayload(payload)
		}
	}()

	fmt.Printf("waiting for message [Exchange, queue]: [%v, %s]", Exchange, queue.Name)
	<-forever
	return nil
}

func handlePayload(payload *Payload) {
	switch payload.Name {
	case "log", "event":
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		//TODO
		fallthrough
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(payload *Payload) error {
	const (
		logServiceURL = "http://logger-service/log"
	)
	buf, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	newReq, err := http.NewRequest(http.MethodPost, logServiceURL, bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	newReq.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: time.Second * 180,
	}

	res, err := client.Do(newReq)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return err
	}

	return nil

}

func NewConsummer(conn *amqp.Connection) (Consummer, error) {
	consumer := Consummer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consummer{}, err
	}

	return consumer, nil
}
