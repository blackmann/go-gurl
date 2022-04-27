package lib

import (
	"net/http"
)

type Request struct {
	Address
	Headers http.Header
	Body    string
}
