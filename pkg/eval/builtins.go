package eval

import "github.com/fugiman/tyrantbot/pkg/message"

type evalFunc func(...interface{}) interface{}

var builtins = map[string]evalFunc{
	"msg":            nil, // Stub for the passed in message
	"msg_RawContent": func(args ...interface{}) interface{} { return args[0].(*message.Message).RawContent },
	"msg_Login":      func(args ...interface{}) interface{} { return args[0].(*message.Message).Login },

	"len":    _len,
	"append": _append,
	"get":    get,
	"set":    set,
	"send":   send,
}

func _len(args ...interface{}) interface{} {
	return 0 // len(args[0].(builtin.Type))
}

func _append(args ...interface{}) interface{} {
	return args[0]
}

func get(args ...interface{}) interface{} {
	return []string{}
}

func set(args ...interface{}) interface{} {
	return nil
}

func send(args ...interface{}) interface{} {
	//log.Printf(args[0].(string), args[1:]...)
	return nil
}
