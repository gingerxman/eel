package handler

import (
	"context"
	"github.com/gingerxman/eel/event"
	"net/http"
)

// NewContext return the Context with Input and Output
func NewContext() *Context {
	return &Context{}
}

type Context struct {
	Request  *Request
	Response *Response
	Data     map[string]interface{}
	Events   []map[string]interface{}
}

// Reset init Context
func (ctx *Context) Reset(rw http.ResponseWriter, r *http.Request) {
	if ctx.Request == nil {
		ctx.Request = &Request{}
	}
	ctx.Request.Reset(r)

	if ctx.Response == nil {
		ctx.Response = &Response{}
	}
	ctx.Response.Reset(rw)

	ctx.Data = nil
	ctx.Events = nil
}

func (ctx *Context) Set(key string, value interface{}) {
	if ctx.Data == nil {
		ctx.Data = make(map[string]interface{})
	}

	ctx.Data[key] = value
}

func (ctx *Context) Get(key string) interface{} {
	return ctx.Data[key]
}

func (ctx *Context) AddEvent(eventData map[string]interface{}) {
	if eventData == nil {
		return
	}
	ctx.Events = append(ctx.Events, eventData)
}

func (ctx *Context) EmitAll() {
	for _, eventData := range ctx.Events {
		event.EmitLocalEvent(eventData)
	}
	ctx.Events = nil
}

func (ctx *Context) SetJSON(key string, value map[string]interface{}) {
	ctx.Request.SetJSON(key, value)
}

func (ctx *Context) SetJSONArray(key string, value []interface{}) {
	ctx.Request.SetJSONArray(key, value)
}

func (ctx *Context) GetBusinessContext() context.Context {
	value := ctx.Get("bContext")
	if value == nil {
		return nil
	} else {
		return value.(context.Context)
	}
}

func (ctx *Context) SetBusinessContext(value interface{}) {
	ctx.Set("bContext", value)
}
