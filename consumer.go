package main

import (
	"github.com/streadway/amqp"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func CreateConsumer(uri, exchange, exchangeType, queueName, key, ctag string) (c *Consumer, err error) {
	c = &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}
	c.conn, err = amqp.Dial(uri)
	if err != nil {
		return
	}
	c.channel, err = c.conn.Channel()
	if err != nil {
		return
	}
	err = c.channel.ExchangeDeclare(
		exchange,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return
	}
	queue, err := c.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return
	}
	err = c.channel.QueueBind(
		queue.Name,
		key,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return
	}

	return
}
