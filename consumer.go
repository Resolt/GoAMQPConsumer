package main

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

type Consumer struct {
	ac           *amqp.Connection
	uri          string
	queueName    string
	exchangeName string
	done         chan error
}

func getConsumer() (c *Consumer, err error) {
	host := os.Getenv("AMQP_HOST")
	user := os.Getenv("AMQP_USER")
	pass := os.Getenv("AMQP_PASS")
	port := os.Getenv("AMQP_PORT")
	vhost := os.Getenv("AMQP_VHOST")

	c = &Consumer{
		uri:          fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, pass, host, port, vhost),
		queueName:    os.Getenv("AMQP_QUEUE"),
		exchangeName: os.Getenv("AMQP_EXCHANGE"),
		done:         make(chan error),
	}
	c.ac, err = amqp.Dial(c.uri)
	if err != nil {
		return
	}
	ch, err := c.ac.Channel()
	if err != nil {
		return
	}
	defer func() { err = ch.Close() }()

	err = ch.ExchangeDeclare(
		c.exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return
	}

	_, err = ch.QueueDeclare(
		c.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return
	}

	return
}
