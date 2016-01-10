package main

import (
	"encoding/json"
	"fmt"

	"github.com/donovanhide/eventsource"
)

type Ingest interface {
	Messages() <-chan *Message
	Errors() <-chan error
}

type firehoseIngest struct {
	es     *eventsource.Stream
	output chan *Message
	errors chan error
}

func NewFirehoseIngest(login string, token string) (Ingest, error) {
	url := fmt.Sprintf("http://tmi.twitch.tv/firehose?login=%s&oauth_token=%s", login, token)
	s, err := eventsource.Subscribe(url, "")
	if err != nil {
		return nil, err
	}

	f := &firehoseIngest{
		es:     s,
		output: make(chan *Message, 1024),
		errors: make(chan error, 1024),
	}
	f.run()
	return f, nil
}

func (f *firehoseIngest) Messages() <-chan *Message {
	return f.output
}

func (f *firehoseIngest) Errors() <-chan error {
	return f.errors
}

func (f *firehoseIngest) run() {
	go func() {
		for e := range f.es.Events {
			switch e.Event() {
			case "privmsg":
				m := &Message{}
				err := json.Unmarshal([]byte(e.Data()), m)
				if err == nil {
					m.ParseTags()
					m.ParseMessage()
					if m.Login != config.Username {
						f.output <- m
					}
				} else {
					f.errors <- err
				}
			default:
				f.errors <- fmt.Errorf("Unknown event type: %v", e)
			}
		}
	}()
	go func() {
		for e := range f.es.Errors {
			f.errors <- e
		}
	}()
}
