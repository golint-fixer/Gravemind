package main

// TODO:
// - Remove oauth token!!! Put it somewhere secret.
// - Actually set Message.IsAction and strip \x01ACTION

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

var config Config

func init() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Create the ingest
	ingest, err := NewFirehoseIngest(config.Username, config.Token)
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
	outgest.Connect(config.Username, config.Token)

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
