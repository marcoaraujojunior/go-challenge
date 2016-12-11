package invoice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/marcoaraujojunior/go-challenge/model"
)

func ListInvoices(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(model.GetAll())
}

func GetInvoice(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "GetInvoice!")
}

func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "CreateInvoice!")
}

func UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "UpdateInvoice!")
}

func DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "DeleteInvoice!")
}

