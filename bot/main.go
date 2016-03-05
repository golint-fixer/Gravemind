package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
)

const testRoom = "#fugitest"

type Config struct {
	IngestName   string `json:"ingest_name"`
	IngestToken  string `json:"ingest_token"`
	OutgestName  string `json:"outgest_name"`
	OutgestToken string `json:"outgest_token"`
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
	// Ayyy pprof!
	go func() {
		http.ListenAndServe(":9999", nil)
	}()

	// Create the ingest
	ingest, err := NewFirehoseIngest(config.IngestName, config.IngestToken)
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
	outgest.Connect(config.OutgestName, config.OutgestToken)

	// Create the brain
	brain, err := NewBrain(ingest.Messages(), outgest.Send)
	if err != nil {
		log.Fatal("NewBrain: ", err)
	}

	// Handle signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		ingest.Stop()
		signal.Stop(c) // In case of catastrophic failure, allow CTRL+C to work again
	}()

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
