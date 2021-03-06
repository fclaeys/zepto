package web

import (
	"context"
	"encoding/json"
	"github.com/go-zepto/zepto/broker"
	"github.com/go-zepto/zepto/logger"
	"github.com/go-zepto/zepto/web/renderer"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
	// Enable webpack asset feature
	_ "github.com/go-webpack/pongo2"
)

type Context interface {
	context.Context
	Request() *http.Request
	Response() http.ResponseWriter
	Params() map[string]string
	Set(string, interface{})
	SetStatus(status int) Context
	Render(template string) error
	RenderJson(data interface{}) error
	Redirect(url string) error
	Logger() logger.Logger
	Broker() *broker.Broker
	Cookies() *Cookies
	Session() *Session
}

type DefaultContext struct {
	logger logger.Logger
	broker *broker.Broker
	context.Context
	res        http.ResponseWriter
	req        *http.Request
	status     int
	data       *sync.Map
	tmplEngine renderer.Engine
	cookies    *Cookies
	session    *Session
}

func NewDefaultContext() *DefaultContext {
	return &DefaultContext{
		Context: context.Background(),
		data:    &sync.Map{},
		status:  200,
	}
}

// Request is the http request
func (d *DefaultContext) Request() *http.Request {
	return d.req
}

// Response is the http response writer
func (d *DefaultContext) Response() http.ResponseWriter {
	return d.res
}

// Set a value to context. The values defined here are accessible in the template
func (d *DefaultContext) Set(key string, value interface{}) {
	d.data.Store(key, value)
}

// Value returns a value from context
func (d *DefaultContext) Value(key interface{}) interface{} {
	if k, ok := key.(string); ok {
		if v, ok := d.data.Load(k); ok {
			return v
		}
	}
	return d.Context.Value(key)
}

// SetStatus set a http status code before render
func (d *DefaultContext) SetStatus(s int) Context {
	d.status = s
	return d
}

// Render a template
func (d *DefaultContext) Render(template string) error {
	return d.tmplEngine.Render(d.res, d.status, template, d.data)
}

// Render a json
func (d *DefaultContext) RenderJson(data interface{}) error {
	d.res.Header().Set("Content-Type", "application/json")
	d.res.WriteHeader(d.status)
	return json.NewEncoder(d.res).Encode(data)
}

// Redirect to url
func (d *DefaultContext) Redirect(url string) error {
	http.Redirect(d.res, d.req, url, 302)
	return nil
}

// Logger is the logger instance from zepto
func (d *DefaultContext) Logger() logger.Logger {
	return d.logger
}

// Broker is the broker instance from zepto
func (d *DefaultContext) Broker() *broker.Broker {
	return d.broker
}

// Retrieve a map of URL parameters
func (d *DefaultContext) Params() map[string]string {
	return mux.Vars(d.req)
}

// Retrieve request session instance
func (d *DefaultContext) Cookies() *Cookies {
	return d.cookies
}

// Retrieve request session instance
func (d *DefaultContext) Session() *Session {
	return d.session
}
