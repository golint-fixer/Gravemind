package main

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

type Outgest interface {
	Connect(string, string)
	Send(string, string, string)
}

type outgest struct {
	users map[string]*User
}

type User struct {
	name  string
	token string
	ch    chan string
	conns *int64
}

func (u *User) Connect() {
	atomic.AddInt64(u.conns, 1)
	NewConn(u.name, u.token, u.ch)
	atomic.AddInt64(u.conns, -1)
}

func NewOutgest() Outgest {
	o := &outgest{
		users: map[string]*User{},
	}
	go o.run()
	return o
}

func (o *outgest) Connect(username, token string) {
	o.users[username] = &User{
		name:  username,
		token: token,
		ch:    make(chan string, 1000),
		conns: new(int64),
	}
	go o.users[username].Connect()
}

func (o *outgest) Send(username, channel, message string) {
	action := "PRIVMSG"
	if channel[0] != '#' {
		action = "WHISPER"
	}

	// For debugging
	channel = "#fugitest"
	log.Printf(":%s %s %s :%s", username, action, channel, message)

	o.users[username].ch <- fmt.Sprintf("%s %s :%s", action, channel, message)
}

func (o *outgest) run() {
	for {
		for _, u := range o.users {
			if int64(len(u.ch)) > *u.conns {
				go u.Connect()
				time.Sleep(10 * time.Millisecond)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
