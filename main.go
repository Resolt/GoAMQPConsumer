package main

import (
	"fmt"
	"log"
)

func main() {
	c, err := getConsumer()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c.uri)
}
