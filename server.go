package main

import (
	"github.com/valyala/fasthttp"
)

type Server struct {
	fsHandler fasthttp.RequestHandler
}

func (s *Server) HandleRequest(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	switch string(ctx.Path()) {
	default:
		s.fsHandler(ctx)
	}
}
