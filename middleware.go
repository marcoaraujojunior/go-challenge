package main

import (
	"net/http"
)

type middleware struct {
	handler      http.Handler
	handlers []http.HandlerFunc
}

func MiddlewareChain() *middleware {
	return &middleware{handlers: make([]http.HandlerFunc, 0, 0)}
}

func (m *middleware) AddHandlerFunc(h ...http.HandlerFunc) {
	m.handlers = append(m.handlers, h...)
}

func (m *middleware) AddHandler(h http.Handler) {
	m.handler = h
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, h := range m.handlers {
		h.ServeHTTP(w, r)
	}
	m.handler.ServeHTTP(w, r)
}

