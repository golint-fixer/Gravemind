package main

// TODO:
// - Remove oauth token!!! Put it somewhere secret.
// - Actually set Message.IsAction and strip \x01ACTION

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/donovanhide/eventsource"
	"github.com/fugiman/tyrantbot/pkg/message"
)

type MessageEvent struct {
	id   string
	data string
}

func (m *MessageEvent) Id() string    { return m.id }
func (m *MessageEvent) Event() string { return "" }
func (m *MessageEvent) Data() string  { return m.data }

func NewMessageEvent(m *message.Message) *MessageEvent {
	return &MessageEvent{
		id:   fmt.Sprint(time.Now().UnixNano()),
		data: fmt.Sprint(m),
	}
}

func main() {
	s := eventsource.NewServer()
	defer s.Close()

	f, err := NewFirehoseIngest("tyrantbot", "")
	if err != nil {
		log.Fatal("NewFirehoseIngest: ", err)
	}

	go func() {
		for m := range f.Messages() {
			s.Publish([]string{"messages"}, NewMessageEvent(m))
		}
	}()
	go func() {
		for e := range f.Errors() {
			log.Println(e)
		}
	}()

	http.HandleFunc("/", s.Handler("messages"))
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
