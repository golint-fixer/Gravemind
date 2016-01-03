package main

type Brain interface {
	Run() <-chan error
}

type brain struct {
	in    <-chan *Message
	out   func(string, string)
	rules map[string][]uint64
	p     *Pool
}

func NewBrain(in <-chan *Message, out func(string, string)) (Brain, error) {
	p := NewPool(64)

	id, err := p.Add(`
    if (msg.UserType == "staff") {
      say(msg.Username + ': ' + msg.RawContent)
      reply('' + HTML(msg.Content))
    }
  `)
	if err != nil {
		return nil, err
	}

	return &brain{
		in:  in,
		out: out,
		rules: map[string][]uint64{
			"": []uint64{id},
		},
		p: p,
	}, nil
}

func (b *brain) Run() <-chan error {
	errs := make(chan error, 1000)
	go func() {
		for m := range b.in {
			say := func(s string) {
				b.out(m.Room, s)
			}
			reply := func(s string) {
				b.out(m.Login, s)
			}
			run := func(c uint64) {
				err := b.p.Run(c, m, say, reply)
				if err != nil {
					errs <- err
				}
			}
			if rules, ok := b.rules[""]; ok {
				for _, c := range rules {
					go run(c)
				}
			}

			if rules, ok := b.rules[m.Room]; ok {
				for _, c := range rules {
					go run(c)
				}
			}
		}
		close(errs)
	}()
	return errs
}