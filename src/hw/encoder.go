package hw

import (
	"encoding/json"
	"log"
	"net/http"
)

type Encoder func(res any) Resp

func DefaultEncoder(res any) Resp {
	w := NewResponse()
	// default decoder
	encoded, err := json.Marshal(res)
	if err != nil {
		log.Printf("error parsing server response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return w
	}
	w.WriteHeader(http.StatusOK)
	w.Write(encoded)
	w.Header().Set("Content-Type", "application/json")
	return w
}
