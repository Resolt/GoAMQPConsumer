package main

import "github.com/sirupsen/logrus"

func main() {
	log := logrus.StandardLogger()

	c, err := getConsumer()
	if err != nil {
		log.Fatal(err)
	}
	c.log = log

	err = c.run()
	if err != nil {
		c.log.Fatal(err)
	}
}
