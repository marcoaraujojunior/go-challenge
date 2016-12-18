package main

import(
	"log"
	"net/http"
	"services/database"
	"services/route"
	"model"
)

func init() {
	database.GetDb().AutoMigrate(&model.Invoice{})
}

func main() {

	router := route.NewRouter()

	log.Fatal(http.ListenAndServe(":80", router))
}

