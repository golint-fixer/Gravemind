package eval

import (
	"github.com/fugiman/tyrantbot/pkg/message"
)

type Expr func(*message.Message, map[string]interface{}) interface{}
type Code []Expr

func (code *Code) Run(msg *message.Message) {
	vars := make(map[string]interface{})
	for _, f := range *code {
		f(msg, vars)
	}
}
