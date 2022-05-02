package lib

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client interface {
	MakeRequest(request Request) (Response, error)
}

type DefaultClient struct {
	client *http.Client
}

func NewHttpClient() DefaultClient {
	return DefaultClient{client: &http.Client{}}
}

func (c DefaultClient) MakeRequest(request Request) (Response, error) {
	start := time.Now()

	req, err := http.NewRequest(request.Address.Method,
		request.Address.Url,
		strings.NewReader(request.Body))

	req.Header = request.Headers

	res, err := c.client.Do(req)

	if err != nil {
		log.Println("Error occurred", err)
		return Response{}, err
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)

	if err != nil {
		log.Println("Error occurred while reading response body", err)
		return Response{}, nil
	}

	timeTaken := time.Since(start).Milliseconds()

	return Response{
		Body:    responseBody,
		Headers: res.Header,
		Status:  res.StatusCode,
		Time:    timeTaken,
		Request: request,
	}, nil
}
