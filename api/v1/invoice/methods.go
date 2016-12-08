package invoice

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Index!")
}

func ListInvoices(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("{test: test}")
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

