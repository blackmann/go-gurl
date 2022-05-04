package lib

import (
	"github.com/blackmann/go-gurl/fork/http/cookiejar"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client interface {
	MakeRequest(request Request) (Response, error)
	GetRawCookies() string
}

type DefaultClient struct {
	client *http.Client
	jar    *cookiejar.Jar
}

func NewHttpClient(rawCookies string) DefaultClient {
	jar := cookiejar.Jar{}
	jar.LoadCookies(strings.NewReader(rawCookies))

	return DefaultClient{client: &http.Client{Jar: &jar}, jar: &jar}
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
		Cookies: res.Cookies(),
		Headers: res.Header,
		Status:  res.StatusCode,
		Time:    timeTaken,
		Request: request,
	}, nil
}

func (c DefaultClient) GetRawCookies() string {
	res := strings.Builder{}
	c.jar.AllCookies(&res)

	return res.String()
}
