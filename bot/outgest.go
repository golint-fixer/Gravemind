package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Outgest interface {
	Connect(string, string)
	Send(string, string) func(string)
}

type outgest struct {
	users   map[string]*User
	pending chan func()
}

type User struct {
	name     string
	token    string
	messages *ConnPool
	whispers *ConnPool
}

type ConnPool struct {
	outgest *outgest
	user    *User
	action  string
	addr    string
	ch      chan string
	count   *int64
}

func NewConnPool(outgest *outgest, user *User, action, addr string) *ConnPool {
	return &ConnPool{
		outgest: outgest,
		user:    user,
		action:  action,
		addr:    addr,
		ch:      make(chan string, 1000),
		count:   new(int64),
	}
}

func (c *ConnPool) Send(channel, message string) {
	c.ch <- fmt.Sprintf("%s %s :%s", c.action, channel, message)
	if int64(len(c.ch)) > *c.count {
		atomic.AddInt64(c.count, 1)
		c.outgest.pending <- func() {
			NewConn(c.addr, c.user.name, c.user.token, c.ch)
			atomic.AddInt64(c.count, -1)
		}
	}
}

func NewOutgest() Outgest {
	o := &outgest{
		users:   map[string]*User{},
		pending: make(chan func(), 1000),
	}
	go o.run()
	return o
}

func (o *outgest) Connect(username, token string) {
	u := &User{
		name:  username,
		token: token,
	}
	u.messages = NewConnPool(o, u, "PRIVMSG", "irc.twitch.tv:6667")
	u.whispers = NewConnPool(o, u, "WHISPER", "192.16.64.180:443")
	o.users[username] = u
}

func (o *outgest) Send(username, channel string) func(string) {
	u := o.users[username]
	c := u.messages
	if channel[0] != '#' {
		c = u.whispers
	}

	return func(message string) {
		c.Send(channel, message)
	}
}

func (o *outgest) run() {
	for fn := range o.pending {
		go fn()
		time.Sleep(1 * time.Second)
	}
}
