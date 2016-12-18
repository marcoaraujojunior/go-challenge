package invoice

import (
	"encoding/json"
	"net/http"
	"net/url"
	"model"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"strings"
	"strconv"
)

func buildLink(q url.Values, pp string, p string, rel string, w http.ResponseWriter, r *http.Request) string {
	var params []string
	q["page"] = append([]string{}, p)
	q["per_page"] = append([]string{}, pp)
	for key, val := range q {
		params = append(params, key + "=" + val[0])
	}
	return "<" + r.Host + r.URL.Path + "?" +strings.Join(params, "&") + "> ; rel=\"" + rel + "\""
}

func buildLinkNext(q url.Values, pp string, p string, t string, w http.ResponseWriter, r *http.Request) string {
	page, _:= strconv.Atoi(p)
	total, _ := strconv.Atoi(t)
	nextPage := page + 1
	if (nextPage > total) {
		return ""
	}
	return buildLink(q, pp, strconv.Itoa(nextPage), "next", w, r)
}

func buildLinkPrev(q url.Values, pp string, p string, w http.ResponseWriter, r *http.Request) string {
	page, _:= strconv.Atoi(p)
	prevPage := page - 1
	if (prevPage < 1) {
		return ""
	}
	return buildLink(q, pp, strconv.Itoa(prevPage), "prev", w, r)
}

func showLinks(invoices model.Invoices, q url.Values, w http.ResponseWriter, r *http.Request) {
	var link []string
	var lf, ll, ln, lp, pp, p, t string

	pp = strconv.Itoa(invoices.PerPage)
	p = strconv.Itoa(invoices.Page)
	t = strconv.Itoa(invoices.Total)

	if invoices.Page > 1 {
		lf = buildLink(q, pp, "1", "first", w, r)
	}

	if invoices.Page < invoices.Total {
		ll = buildLink(q, pp, t, "last", w, r)
	}

	ln = buildLinkNext(q, pp, p, t, w, r)
	lp = buildLinkPrev(q, pp, p, w, r)

	if ln != "" {
		link = append(link, ln)
	}

	if ll != "" {
		link = append(link, ll)
	}

	if lf != "" {
		link = append(link, lf)
	}

	if lp != "" {
		link = append(link, lp)
	}

	if strings.Join(link, ",") != "" {
		w.Header().Set("Link", strings.Join(link, ","))
	}
}

func ListInvoices(w http.ResponseWriter, r *http.Request) {
	var invoices model.Invoices
	q := r.URL.Query()
	invoices = model.GetAll(q)
	showLinks(invoices, q, w, r)
	json.NewEncoder(w).Encode(invoices.Records)
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
	var invoice model.Invoice
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &invoice); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	err = model.Create(&invoice)
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	var invoice model.Invoice
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &invoice); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	invoiceNumber := mux.Vars(r)["invoiceNumber"]
	invoice.Document = invoiceNumber
	err = model.Update(&invoice)
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

