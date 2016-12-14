package invoice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"model"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
)

func ListInvoices(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(model.GetAll())
}

func GetInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceNumber := mux.Vars(r)["invoiceNumber"]
	invoice, err := model.Get(invoiceNumber)
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoice)
}

func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "CreateInvoice!")
}

func UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	var invoice model.Invoice
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &invoice); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	err = model.Save(invoice)
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceNumber := mux.Vars(r)["invoiceNumber"]
	err := model.Delete(invoiceNumber)
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

