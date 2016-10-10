package server

import (
	"github.com/kimiazhu/log4go"
	"github.com/streadway/amqp"
)

func CreateRabbitQueue(url, name string) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log4go.Error("Failed to connect to RabbitMQ %s: ", url, err)
		log4go.Close()
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log4go.Error("Failed to open a channel: ", err)
		panic(err)

	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log4go.Error("Failed to open a channel: ", err)
		panic(err)
	}
}
