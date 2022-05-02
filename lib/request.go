package lib

import (
	"net/http"
)

type Request struct {
	Address
	Headers http.Header
	Body    string
}

type RequestError struct {
	Err     error
	Request *Request // ref to the original request
}
