package main

import(
	"log"
	"net/http"
)

func main() {

	connect()

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":80", router))
}
