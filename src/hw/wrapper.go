package hw

import "net/http"

type Wrapper[TReq any] func(args TReq, r *http.Request) any
