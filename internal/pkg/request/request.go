package request

import (
	"net/http"
)

var Request *request = new(request)

type second int64

const (
	GET_REQUEST    string = "GET"
	POST_REQUEST   string = "POST"
	DELETE_REQUEST string = "DELETE"
)

var IsSkipVerifyDefault bool = true
var Timeout second = 20

type request struct {
	Host     string
	Port     string
	Method   string
	IsVerify bool
	Headers  map[string]string
}

func NewRequest(method, token, url string) (*http.Request, error) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/yaml"
	headers["Authorization"] = "Bearer " + token

	r, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		r.Header.Add(key, value)
	}

	return r, nil
}
