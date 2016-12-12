package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"api/v1/invoice"
)

var routes = invoice.Get()

func NewRouter() *middleware {
	r := loadRoutes()
	m := MiddlewareChain()
	m.AddHandlerFunc(token)
	m.AddHandler(r)

	return m
}

func loadRoutes() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			PathPrefix(route.PathPrefix).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

func token(rw http.ResponseWriter, r *http.Request) {
	log.Println("running foo.")
}

