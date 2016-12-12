package main

import(
	"log"
	"net/http"
	"services/database"
	"model"
)

func init() {
	database.Connect()
	database.Db.AutoMigrate(&model.Invoice{})
}

func main() {

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":80", router))

//	router := NewRouter()
//	http.Handle("/", Middleware(router))

//	log.Fatal(http.ListenAndServe(":80", router))
}


