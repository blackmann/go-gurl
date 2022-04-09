package handler

type Status int

var (
	IDLE       Status = 0
	PROCESSING Status = 1
)

type RequestHandler struct {
	Status Status
}

func NewRequestHandler() RequestHandler {
	return RequestHandler{}
}

func (handler *RequestHandler) MakeRequest() {
	handler.Status = PROCESSING
}
