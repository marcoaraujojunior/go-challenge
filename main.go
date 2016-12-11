package main

import(
	"log"
	"net/http"
	"github.com/marcoaraujojunior/go-challenge/database"
	"github.com/marcoaraujojunior/go-challenge/model"
)

func init() {
	database.Connect()
	database.Db.AutoMigrate(&model.Invoice{})
}

func main() {
	router := NewRouter()

	log.Fatal(http.ListenAndServe(":80", router))
}
