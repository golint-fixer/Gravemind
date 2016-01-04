package main

// TODO:
// - Remove oauth token!!! Put it somewhere secret.
// - Actually set Message.IsAction and strip \x01ACTION

import (
	"log"
	"os"
)

var username = "tyrantbot"
var token = os.Getenv("TOKEN")

func main() {
	// Create the ingest
	ingest, err := NewFirehoseIngest(username, token)
	if err != nil {
		log.Fatal("NewFirehoseIngest: ", err)
	}
	go func() {
		for e := range ingest.Errors() {
			log.Println(e)
		}
	}()

	// Create the outgest
	outgest := NewOutgest()
	outgest.Connect(username, token)

	// Create the brain
	brain, err := NewBrain(ingest.Messages(), outgest.Send)
	if err != nil {
		log.Fatal("NewBrain: ", err)
	}

	// And... go
	for err := range brain.Run() {
		log.Println(err)
	}
}

func logErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
