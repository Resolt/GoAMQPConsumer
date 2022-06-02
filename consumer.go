package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type consumer struct {
	ac           *amqp.Connection
	ch           *amqp.Channel
	uri          string
	queueName    string
	exchangeName string
	tag          string
	done         chan error
	deliveries   <-chan amqp.Delivery
	ntasks       int
	log          *logrus.Logger
}

func getConsumer() (c *consumer, err error) {
	host := os.Getenv("AMQP_HOST")
	user := os.Getenv("AMQP_USER")
	pass := os.Getenv("AMQP_PASS")
	port := os.Getenv("AMQP_PORT")
	vhost := os.Getenv("AMQP_VHOST")

	c = &consumer{
		uri:          fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, pass, host, port, vhost),
		queueName:    os.Getenv("AMQP_QUEUE"),
		exchangeName: os.Getenv("AMQP_EXCHANGE"),
		tag:          os.Getenv("TAG"),
		done:         make(chan error),
	}
	c.ac, err = amqp.Dial(c.uri)
	if err != nil {
		return
	}
	c.ch, err = c.ac.Channel()
	if err != nil {
		return
	}

	ntasks, err := strconv.Atoi(os.Getenv("NTASKS"))
	if err != nil {
		return
	}
	if ntasks <= 0 {
		c.ntasks = runtime.NumCPU()
	} else {
		c.ntasks = ntasks
	}

	err = c.ch.Qos(c.ntasks*2, 0, false)
	if err != nil {
		return
	}

	err = c.ch.ExchangeDeclare(
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

	_, err = c.ch.QueueDeclare(
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

	c.deliveries, err = c.ch.Consume(
		c.queueName,
		c.tag,
		false,
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

func (c *consumer) run() (err error) {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go c.handle(ctx)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case <-interrupt:
		c.log.Info("sigterm received")
		err = c.shutdown()
		<-c.done
	case <-c.done:
		err = c.shutdown()
	}

	return
}

func (c *consumer) handle(ctx context.Context) {
	wg := sync.WaitGroup{}
	guard := make(chan struct{}, int(c.ntasks))
	for d := range c.deliveries {
		select {
		case <-ctx.Done():
			break
		default:
			guard <- struct{}{}
			wg.Add(1)
			go func(d amqp.Delivery) {
				// c.log.Info("Message Body: ", string(d.Body))
				c.log.Info(d.Body)
				d.Ack(false)
				time.Sleep(time.Second)
				wg.Done()
				<-guard
			}(d)
		}
	}
	wg.Wait()
	c.done <- nil
}

func (c *consumer) shutdown() (err error) {
	c.log.Info("shutting down")
	err = c.ch.Cancel(c.tag, false)
	if err != nil {
		return
	}
	err = c.ch.Close()
	if err != nil {
		return
	}
	err = c.ac.Close()
	if err != nil {
		return
	}
	return
}
