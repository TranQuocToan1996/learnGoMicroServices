package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// an exchange starting with "amq."
// Exchanges are message routing agents, defined by the virtual host within RabbitMQ. An exchange is responsible for routing the messages to different queues with the help of header attributes, bindings, and routing keys. A binding is a "link" that you set up to bind a queue to an exchange
func declareExchange(ch *amqp.Channel) error {

	return ch.ExchangeDeclare(
		Exchange, // name
		"topic",  // Type
		true,     // durable -> like keep-alive
		false,    // autoDelete
		false,    // UseInternal -> no, we using in our microservices
		false,    // nowait -> will assume to be declared on the server
		nil,      // arg table (extra parameters)
	)
	/* Durable and Non-Auto-Deleted exchanges will survive server restarts and remain
	   declared when there are no remaining bindings.  This is the best lifetime for
	   long-lived exchange configurations like stable routes and default exchanges.

	   Non-Durable and Auto-Deleted exchanges will be deleted when there are no
	   remaining bindings and not restored on server restart.  This lifetime is
	   useful for temporary topologies that should not pollute the virtual host on
	   failure or after the consumers have completed.

	   Non-Durable and Non-Auto-deleted exchanges will remain as long as the server is
	   running including when there are no remaining bindings.  This is useful for
	   temporary topologies that may have long delays between bindings.

	   Durable and Auto-Deleted exchanges will survive server restarts and will be
	   removed before and after server restarts when there are no remaining bindings.
	   These exchanges are useful for robust temporary topologies or when you require
	   binding durable queues to auto-deleted exchanges. */
}

// declare queue to hold messages and deliver to consumers
func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name
		false, // durable -> get rid of it when unuse
		false, // autoDelete -> delete when unuse
		true,  // Exclusive -> do not share around
		false, // nowait -> will assume to be declared on the server
		nil,   // arg table (extra parameters)
	)
}
