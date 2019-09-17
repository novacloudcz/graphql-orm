package main

import (
	"log"

	"github.com/novacloudcz/graphql-orm/cmd"
	"github.com/novacloudcz/graphql-orm/events"
)

func main() {
	cmd.Execute()
}

// this is just for importing the events package and adding it to the go modules
func testEventController() {
	_, err := events.NewEventController()
	log.Println(err)
}
