package request

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

var Request *request = new(request)

type method string
type second int64

const (
	GET_REQUEST    method = "GET"
	POST_REQUEST   method = "POST"
	DELETE_REQUEST method = "DELETE"

	TIMEOUT_REQUEST second = 10
)

type request struct {
	Host     string
	Port     string
	Method   method
	IsVerify bool
	Headers  map[string]string
}

func NewRequest(req *request, url string) (*http.Request, error) {
	r, err := http.NewRequest(string(req.Method), url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range req.Headers {
		r.Header.Add(key, value)
	}

	return r, nil
}

func HandleRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !Request.IsVerify,
			},
		},
		Timeout: time.Duration(TIMEOUT_REQUEST) * time.Second,
	}
	defer client.CloseIdleConnections()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func Handler(url string) ([]byte, error) {
	req, err := NewRequest(Request, url)
	if err != nil {
		return nil, err
	}

	resp, err := HandleRequest(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func InitLogRequest(host, port, token string, isVerify bool) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/yaml"
	headers["Authorization"] = "Bearer " + token

	Request.Host = host
	Request.Port = port
	Request.IsVerify = isVerify
	Request.Headers = headers
	Request.Method = GET_REQUEST
}
