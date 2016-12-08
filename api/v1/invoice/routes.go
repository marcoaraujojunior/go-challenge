package invoice

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func Get() Routes {
	return routes
}

var routes = Routes {
	Route{
		"Index",
		"GET",
		"/v1",
		Index,
	},
	Route{
		"ListInvoices",
		"GET",
		"/v1/invoices",
		ListInvoices,
	},
	Route{
		"GetInvoice",
		"GET",
		"/v1/invoice/{invoiceNumber}",
		GetInvoice,
	},
	Route{
		"CreateInvoice",
		"POST",
		"/v1/invoice",
		CreateInvoice,
	},
	Route{
		"UpdateInvoice",
		"PUT",
		"/v1/invoice/{invoiceNumber}",
		UpdateInvoice,
	},
	Route{
		"UpdateInvoice",
		"DELETE",
		"/v1/invoice/{invoiceNumber}",
		DeleteInvoice,
	},
}

