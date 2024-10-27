package hw

import (
	"net/http"
)

type resp struct {
	header     http.Header
	body       []byte
	statusCode int
}

type Resp interface {
	Header() http.Header
	Write(b []byte) (int, error)
	WriteHeader(statusCode int)
	WriteResponse(w http.ResponseWriter)
}

func (resp *resp) Header() http.Header {
	return resp.header
}

func (resp *resp) Write(b []byte) (int, error) {
	resp.body = b
	return len(b), nil
}

func (resp *resp) WriteHeader(statusCode int) {
	resp.statusCode = statusCode
}

func (resp *resp) WriteResponse(w http.ResponseWriter) {
	for key, values := range resp.header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.statusCode)
	if resp.body != nil {
		w.Write(resp.body)
	}
}

func NewResponse() Resp {
	return &resp{
		statusCode: http.StatusOK,
		body:       []byte{},
		header:     http.Header{},
	}
}
