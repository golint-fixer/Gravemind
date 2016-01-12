package dynamosync

import (
	"reflect"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Cache struct {
	template reflect.Type
	records  Syncronizer
	data     map[interface{}]interface{}
}

func NewCache(sess *session.Session, table string, structValue interface{}) (*Cache, error) {
	structType := reflect.TypeOf(structValue)
	records, err := New(sess, table)
	if err != nil {
		return nil, err
	}
	c := &Cache{
		template: structType,
		records:  records,
		data:     map[interface{}]interface{}{},
	}
	go c.run()
	return c, nil
}

func (c *Cache) Get(key interface{}) interface{} {
	return c.data[key]
}

func (c *Cache) run() {
	var keys []string
	for i := 0; i < c.template.NumField(); i++ {
		field := c.template.Field(i)
		if field.Tag.Get("dynamosync") == "key" {
			keys = append(keys, field.Name)
		}
	}
	for r := range c.records {
		rk := reflect.New(c.template).Elem()
		rv := reflect.New(c.template)
		v := rv.Interface()
		dynamodbattribute.ConvertFromMap(r, v)
		for _, name := range keys {
			rk.FieldByName(name).Set(rv.Elem().FieldByName(name))
		}
		c.data[rk.Interface()] = v
	}
}
