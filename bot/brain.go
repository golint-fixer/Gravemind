package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fugiman/tyrantbot/shared/dynamosync"
)

type Brain interface {
	Run() <-chan error
}

type codeKey struct {
	id        int
	createdAt int
}

type brain struct {
	sync.RWMutex

	in       <-chan *Message
	out      func(string, string) func(string)
	channels map[string][]uint64
	codes    map[codeKey]uint64
	p        *Pool
	plugins  dynamosync.Syncronizer
	users    dynamosync.Syncronizer
}

var bots = map[string]string{}

func NewBrain(in <-chan *Message, out func(string, string) func(string)) (Brain, error) {
	sess := session.New(aws.NewConfig().WithRegion("us-west-2"))

	plugins, err := dynamosync.New(sess, "plugins")
	if err != nil {
		return nil, err
	}

	users, err := dynamosync.New(sess, "users")
	if err != nil {
		return nil, err
	}

	return &brain{
		in:       in,
		out:      out,
		channels: map[string][]uint64{},
		codes:    map[codeKey]uint64{},
		p:        NewPool(64),
		plugins:  plugins,
		users:    users,
	}, nil
}

func (b *brain) Run() <-chan error {
	errs := make(chan error, 1000)
	addError := func(err error) {
		if err != nil {
			errs <- err
		}
	}
	go func() {
	SeedLoop:
		for {
			select {
			case p := <-b.plugins:
				addError(b.addPlugin(p))
			default:
				break SeedLoop
			}
		}
	RunLoop:
		for {
			select {
			case m, ok := <-b.in:
				if !ok {
					break RunLoop
				}
				go b.handleMessage(m, errs)
			case p := <-b.plugins:
				addError(b.addPlugin(p))
			case u := <-b.users:
				addError(b.addUser(u))
			}
		}
		time.Sleep(3 * time.Second)
		close(errs)
	}()
	return errs
}

type plugin struct {
	UserId    int    `json:"user_id"`
	CreatedAt int    `json:"created_at"`
	Name      string `json:"name"`
	Version   int    `json:"version"`
	Code      string `json:"code"`
}

func (b *brain) addPlugin(data map[string]*dynamodb.AttributeValue) error {
	item := &plugin{}
	err := dynamodbattribute.ConvertFromMap(data, item)
	if err != nil {
		return err
	}
	key := codeKey{item.UserId, item.CreatedAt}
	b.Lock()
	defer b.Unlock()
	id, err := b.p.Add(item.Code)
	if err != nil {
		return err
	}
	b.codes[key] = id
	//log.Printf("b.codes[%v]=%v", key, id)
	return nil
}

type user struct {
	Id      int       `json:"id"`
	Login   string    `json:"login"`
	Name    string    `json:"name"`
	Plugins []*plugin `json:"plugins"`
}

func (b *brain) addUser(data map[string]*dynamodb.AttributeValue) error {
	item := &user{}
	err := dynamodbattribute.ConvertFromMap(data, item)
	if err != nil {
		return err
	}
	key := "#" + item.Login
	b.Lock()
	defer b.Unlock()
	var plugins []uint64
	for _, p := range item.Plugins {
		key := codeKey{p.UserId, p.CreatedAt}
		if id, ok := b.codes[key]; ok {
			plugins = append(plugins, id)
		} else {
			return fmt.Errorf("Couldn't find plugin matching key: %v", key)
		}
	}
	b.channels[key] = plugins
	//log.Printf("b.channels[%v]=%v", key, plugins)
	return nil
}

func (b *brain) handleMessage(m *Message, errs chan error) {
	b.RLock()
	var wg sync.WaitGroup
	bot := config.Username
	if b, ok := bots[m.Room]; ok {
		bot = b
	}
	run := func(p uint64, say, reply func(string)) {
		err := b.p.Run(p, m, say, reply)
		if err != nil {
			errs <- err
		}
		wg.Done()
	}
	if plugins, ok := b.channels["#"]; ok {
		for _, p := range plugins {
			wg.Add(1)
			go run(p, b.out(bot, testRoom), b.out(bot, testRoom))
		}
	}

	if plugins, ok := b.channels[m.Room]; ok {
		for _, p := range plugins {
			wg.Add(1)
			go run(p, b.out(bot, m.Room), b.out(bot, m.Login))
		}
	}
	b.RUnlock()
	wg.Wait()
}
