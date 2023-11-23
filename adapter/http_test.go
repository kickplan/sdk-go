package adapter

import (
	"net/http"
)

type mockClient struct{}

var DoFunc func(req *http.Request) (*http.Response, error)

func (c *mockClient) Do(req *http.Request) (*http.Response, error) {
	return DoFunc(req)
}
