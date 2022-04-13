package handler

import (
	"github.com/blackmann/gurl/common/status"
	"time"
)

type RequestHandler struct {
	onRequestStart func()
	onRequestEnd   func()
	Status         status.Status
}

func NewRequestHandler() RequestHandler {
	return RequestHandler{}
}

func (handler *RequestHandler) MakeRequest() {
	handler.Status = status.PROCESSING

	time.Sleep(3 * time.Second)
	handler.Status = status.IDLE
}

func (handler *RequestHandler) OnStartRequest(f func()) {
	handler.onRequestStart = f
}

func (handler *RequestHandler) OnRequestEnd(f func()) {
	handler.onRequestEnd = f
}
