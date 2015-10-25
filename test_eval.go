package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/fugiman/tyrantbot/pkg/eval"
	"github.com/fugiman/tyrantbot/pkg/message"
)

const ADD_QUOTE = `
quotes := get("quotes")
quotes = append(quotes, msg.RawContent)
set("quotes", quotes)
send("@%s: Added quote #%d", msg.Login, len(quotes))
`

const MESSAGE = `{"command":"","room":"#vgbootcamp","nick":"mayday_believes","body":"BibleThump","tags":"color=#DAA520;display-name=Mayday_Believes;emotes=86:0-9;subscriber=0;turbo=0;user-id=43613324;user-type="}`

func main() {
	m := &message.Message{}
	err := json.Unmarshal([]byte(MESSAGE), m)
	if err != nil {
		log.Fatal(err)
	}

	m.ParseTags()
	m.ParseMessage()

	code, err := eval.Parse(ADD_QUOTE)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	for i := 0; i < 500000; i++ { // 500,000/s
		code.Run(m)
	}
	end := time.Now()

	log.Println(end.Sub(start))
}
