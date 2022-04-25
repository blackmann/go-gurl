package lib

import (
	"io"
	"log"
	"net/http"
	"time"
)

type Client interface {
	MakeRequest(address Address, header http.Header, reader io.Reader) (Response, error)
}

type DefaultClient struct {
	client *http.Client
}

func NewHttpClient() DefaultClient {
	return DefaultClient{client: &http.Client{}}
}

func (c DefaultClient) MakeRequest(address Address, headers http.Header, body io.Reader) (Response, error) {
	start := time.Now()

	req, err := http.NewRequest(address.Method, address.Url, body)
	req.Header = headers

	res, err := c.client.Do(req)
	defer res.Body.Close()

	if err != nil {
		log.Println("Error occurred", err)
		return Response{}, err
	}

	responseBody, err := io.ReadAll(res.Body)

	if err != nil {
		log.Println("Error occurred while reading response body", err)
		return Response{}, nil
	}

	timeTaken := time.Since(start).Milliseconds()

	return Response{Body: responseBody, Headers: res.Header, Status: res.StatusCode, Time: timeTaken}, nil
}
