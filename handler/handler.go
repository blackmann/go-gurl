package handler

import (
	"github.com/blackmann/gurl/common/status"
	"time"
)

type RequestHandler struct {
	Status status.Status
}

func NewRequestHandler() RequestHandler {
	return RequestHandler{}
}

func (handler *RequestHandler) makeRequest() {
	time.Sleep(3 * time.Second)
	handler.Status = status.IDLE
}

func (handler *RequestHandler) MakeRequest() {
	handler.Status = status.PROCESSING

	go handler.makeRequest()
}
