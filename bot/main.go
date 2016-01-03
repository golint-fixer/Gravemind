package main

// TODO:
// - Remove oauth token!!! Put it somewhere secret.
// - Actually set Message.IsAction and strip \x01ACTION

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/donovanhide/eventsource"
)

type Event struct {
	id   string
	data string
}

func (m *Event) Id() string    { return m.id }
func (m *Event) Event() string { return "" }
func (m *Event) Data() string  { return m.data }

func main() {
	// Create the ingest
	ingest, err := NewFirehoseIngest("tyrantbot", os.Getenv("TOKEN"))
	if err != nil {
		log.Fatal("NewFirehoseIngest: ", err)
	}
	go func() {
		for e := range ingest.Errors() {
			log.Println(e)
		}
	}()

	// Create the outgest
	s := eventsource.NewServer()
	defer s.Close()
	http.HandleFunc("/", s.Handler("messages"))
	go func() {
		err = http.ListenAndServe(":8000", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	send := func(channel string, message string) {
		s.Publish([]string{"messages"}, &Event{
			id:   fmt.Sprint(time.Now().UnixNano()),
			data: fmt.Sprintf("<%s> %s", channel, message),
		})
	}

	// Create the brain
	brain, err := NewBrain(ingest.Messages(), send)
	if err != nil {
		log.Fatal("NewBrain: ", err)
	}

	// And... go
	for err := range brain.Run() {
		log.Println(err)
	}
}
