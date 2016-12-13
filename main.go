package main

import(
	"log"
	"net/http"
	"services/database"
	"model"
)

func init() {
	database.GetDb().AutoMigrate(&model.Invoice{})
}

func main() {

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":80", router))
}

