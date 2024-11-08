package hw

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/schema"
)

// any here can be (nil, *resp, T)
type Decoder[T any] func(*http.Request) any // should return T, Resp or nil

func DefaultDecoder[T any](r *http.Request) any {
	var args T
	if r.Method == http.MethodGet {
		decoder := schema.NewDecoder()
		var mapped map[string]any
		err := decoder.Decode(&mapped, r.URL.Query())
		if err == nil {
			encoded, _ := json.Marshal(mapped)
			json.Unmarshal(encoded, &args)
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
