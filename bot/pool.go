package main

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/fugiman/otto"
)

type worker struct {
	sync.Mutex
	vm    *otto.Otto
	codes map[uint64]otto.Value
}

type Pool struct {
	workers   []*worker
	available chan *worker
	counter   *uint64
}

func NewPool(workers int) *Pool {
	p := &Pool{
		workers:   []*worker{},
		available: make(chan *worker, workers),
		counter:   new(uint64),
	}
	for i := 0; i < workers; i++ {
		w := &worker{
			vm:    otto.New(),
			codes: map[uint64]otto.Value{},
		}
		p.workers = append(p.workers, w)
		p.available <- w
	}
	return p
}

func (p *Pool) Add(code string) (uint64, error) {
	id := atomic.AddUint64(p.counter, 1)
	for _, w := range p.workers {
		w.Lock()
		c, err := w.vm.Run(fmt.Sprintf(`x = function(msg, say, reply) {%s}`, code))
		if err == nil {
			w.codes[id] = c
		}
		w.Unlock()
		if err != nil {
			return id, err
		}
	}
	return id, nil
}

func (p *Pool) Run(id uint64, msg *Message, say, reply func(string)) error {
	var err error
	w := <-p.available
	w.Lock()
	if code, ok := w.codes[id]; ok {
		_, err = code.Call(code, msg, say, reply)
	} else {
		err = fmt.Errorf("Invalid code id: %d", id)
	}
	w.Unlock()
	p.available <- w
	return err
}
