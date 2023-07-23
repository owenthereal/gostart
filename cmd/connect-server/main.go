package main

import (
	"net/http"

	"github.com/owenthereal/gostart/proto"
	"github.com/owenthereal/gostart/proto/gen/user/v1/v1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	userSvc := proto.NewUserService()

	mux := http.NewServeMux()
	path, handler := v1connect.NewUserServiceHandler(userSvc)
	mux.Handle(path, handler)
	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
