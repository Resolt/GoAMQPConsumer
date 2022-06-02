package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	lr := logrus.StandardLogger()
	lr.SetFormatter(&logrus.JSONFormatter{})

	c, err := getConsumer()
	if err != nil {
		lr.Fatal(err)
	}
	c.log = lr

	err = c.run()
	if err != nil {
		lr.Fatal(err)
	}
}
