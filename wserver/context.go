package wserver

import (
	"encoding/json"
)

type Context struct {
	Conn *Conn

	storage map[string]interface{}

	Update *Update
}

func (c *Context) Set(key string, val interface{}) {
	c.storage[key] = val
}

func (c *Context) Get(key string) interface{} {
	return c.storage[key]
}

// Converts Context.Update.Data to i with encoding/json
func (c *Context) Bind(i interface{}) error {
	b, err := json.Marshal(c.Update.Data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, i)
}
func (c *Context) Data() interface{} {
	return c.Update.Data
}

func (c *Context) EventType() string {
	return c.Update.EventType
}
