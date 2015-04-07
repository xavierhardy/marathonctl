package main

import (
	"io"
	"net/http"
	"strings"
)

type Action interface {
	// Apply runs the command with args and gets the json result.
	Apply(args []string)
}

type Login struct {
	Host string // todo multi hosts
	User string
	Pass string
}

func NewLogin(host, login string) *Login {
	toks := strings.SplitN(login, ":", 2)
	return &Login{
		Host: host,
		User: toks[0],
		Pass: toks[1],
	}
}

type Tool struct {
	selections map[string]Selector
}

func (t *Tool) Start(args []string) {
	if len(args) == 0 {
		Usage()
	}
	if selection, ok := t.selections[args[0]]; !ok {
		Usage()
	} else {
		selection.Select(args[1:])
	}
}

type Client struct {
	client http.Client
	login  *Login
}

func NewClient(login *Login) *Client {
	return &Client{
		client: http.Client{},
		login:  login,
	}
}

func (c *Client) Do(r *http.Request) (*http.Response, error) {
	return c.client.Do(r)
}

func (c *Client) GET(path string) *http.Request {
	url := c.login.Host + path
	request, e := http.NewRequest("GET", url, nil)
	Check(e == nil, "failed to crete GET request", e)
	c.tweak(request)
	return request
}

func (c *Client) POST(path string, body io.ReadCloser) *http.Request {
	url := c.login.Host + path
	request, e := http.NewRequest("POST", url, body)
	Check(e == nil, "failed to create POST request", e)
	c.tweak(request)
	return request
}

func (c *Client) DELETE(path string) *http.Request {
	url := c.login.Host + path
	request, e := http.NewRequest("DELETE", url, nil)
	Check(e == nil, "failed to create DELETE request", e)
	c.tweak(request)
	return request
}

func (c *Client) PUT(path string, body io.ReadCloser) *http.Request {
	url := c.login.Host + path
	request, e := http.NewRequest("PUT", url, body)
	Check(e == nil, "failed to create PUT request", e)
	c.tweak(request)
	return request
}

// tweak will set:
// Content-Type: application/json
// Accept: application/json
// Accept-Encoding: gzip, deflate, compress
func (c *Client) tweak(request *http.Request) {
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.SetBasicAuth(c.login.User, c.login.Pass)
}
