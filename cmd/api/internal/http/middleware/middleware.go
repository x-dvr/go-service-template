package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

type middlewareChain struct {
	middlewares []Middleware
}

func (mc middlewareChain) For(h http.Handler) http.Handler {
	for i := len(mc.middlewares) - 1; i >= 0; i-- {
		h = mc.middlewares[i](h)
	}

	return h
}

func Use(middlewares ...Middleware) middlewareChain {
	newCons := make([]Middleware, 0, len(middlewares))
	newCons = append(newCons, middlewares...)

	return middlewareChain{newCons}
}
