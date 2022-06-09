package client

import (
	"net/http"
	"net/url"
)

type SimpleHTTPClient interface {
	Get(url string) (resp *http.Response, err error)
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

type MockHTTPClient struct {
	url  string
	data url.Values
	resp *http.Response
	err  error
}

func (mock *MockHTTPClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	mock.url = url
	mock.data = data
	return mock.resp, mock.err
}

func (mock *MockHTTPClient) Get(url string) (resp *http.Response, err error) {
	mock.url = url
	return mock.resp, mock.err
}
