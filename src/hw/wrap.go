package hw

import (
	"log"
	"net/http"
	"reflect"
)

type ParserOptions[T any] struct {
	Decoder Decoder[T]
}

func Wrap[T any](fn Wrapper[T], options *ParserOptions[T]) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var args T
		if reflect.TypeOf(args).Kind() != reflect.Struct {
			panic("hw (http wrap) can parse only structs. tried to handle not a struct")
		}
		// default options values
		if options == nil {
			options = &ParserOptions[T]{Decoder: DefaultDecoder[T]}
		}
		if options.Decoder == nil {
			options.Decoder = DefaultDecoder[T]
		}

		// decoder
		decoded := options.Decoder(r)
		if decoded == nil { // handle decode nil
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		if resp, ok := decoded.(Resp); ok { // handle decode *resp
			resp.WriteResponse(w)
			return
		}

		if arguments, ok := decoded.(T); !ok { // handle decode T
			log.Printf("decoder returns unexpected type returned: %v", arguments)
			w.WriteHeader(500)
			return
		} else {
			args = arguments
		}

		// run endpoint
		res := fn(args, r)

		// encoder
		if res == nil { // nil response
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if resp, ok := res.(Resp); ok { // custom response
			resp.WriteResponse(w)
			return
		}

		// default encode response and respond
		defaultRes := DefaultEncoder(res)
		defaultRes.WriteResponse(w)
	}
}
