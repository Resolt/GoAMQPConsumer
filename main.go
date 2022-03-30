package main

func main() {
	c, err := getConsumer()
	if err != nil {
		logFatal(err)
	}

	err = c.run()
	if err != nil {
		logFatal(err)
	}
}
