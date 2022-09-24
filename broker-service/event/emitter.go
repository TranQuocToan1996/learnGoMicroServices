package event

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer func() {
		closeChanConn(channel)
	}()

	return declareExchange(channel)
}

func closeChanConn(channel *amqp.Channel) {
	var count int = 10
	for count > 0 {
		time.Sleep(time.Second)
		count--
		err := channel.Close()
		if err != nil {
			log.Println(err)
			continue
		}
		return
	}
}

func (e *Emitter) Push(ctx context.Context, event, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer func() {
		closeChanConn(channel)
	}()

	log.Println("Pushing to channel")

	if ctx == nil {
		ctx = context.TODO()
	}

	/* Publishings can be undeliverable when the mandatory flag is true and no queue is
	bound that matches the routing key, or when the immediate flag is true and no
	consumer on the matched queue is ready to accept the delivery. */
	err = channel.PublishWithContext(ctx,
		Exchange, // Exchange
		severity, // routing key
		false,    // mandatory
		false,    // immedialy
		amqp.Publishing{
			ContentType: "text/plain", // string
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
