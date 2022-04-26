package lib

import (
	"io"
	"net/http"
)

type Request struct {
	Address
	Headers http.Header
	Body    io.Reader
}
