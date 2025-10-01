package types

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter
	Params  map[string]string
}

type HandlerFunc func(*Context)

type TokenValidator interface {
	ValidateToken(token string) bool
}

type Middleware func(HandlerFunc) HandlerFunc

func (c *Context) JSON(status int, data interface{}) error {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)
	return json.NewEncoder(c.Writer).Encode(data)
}

func (c *Context) BindJSON(dest interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(dest)
}

func (c *Context) QueryParam(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) FormValue(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) Error(status int, message string) {
	c.JSON(status, map[string]string{"error": message})
}
