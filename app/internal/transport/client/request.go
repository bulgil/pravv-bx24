package client

import "net/http"

type Request struct {
	Request *http.Request
	Resp    chan *http.Response
	Error   chan error
}
