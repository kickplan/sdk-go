package adapter

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type mockClient struct {
	mock.Mock
}

var DoFunc func(req *http.Request) (*http.Response, error)

func (c *mockClient) Do(req *http.Request) (*http.Response, error) {
	return DoFunc(req)
}
