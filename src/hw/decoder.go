package hw

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/schema"
)

// any here can be (nil, *resp, T)
type Decoder[T any] func(*http.Request) any // should return T or nil

func DefaultDecoder[T any](r *http.Request) any {
	var args T
	if r.Method == http.MethodGet {
		decoder := schema.NewDecoder()
		err := decoder.Decode(&args, r.URL.Query())
		if err == nil {
			return args
		}
		res := NewResponse()
		res.WriteHeader(http.StatusBadRequest)
		io.WriteString(res, "invalid search params")
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&args); err != nil {
		return nil
	}
	return args
}
