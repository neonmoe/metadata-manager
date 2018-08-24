package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

func main() {
	var fs = &fasthttp.FS{
		Root:       "./public",
		IndexNames: []string{"index.html"},
	}
	var fsHandler = fs.NewRequestHandler()

	var server = &Server{
		fsHandler: fsHandler,
	}

	var addr = "0.0.0.0:8080"
	fmt.Printf("Serving at: %s", addr)
	fasthttp.ListenAndServe(addr, server.HandleRequest)
}
