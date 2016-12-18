package main_test

import(
	"log"
	"bytes"
	"testing"
	"net/url"
	"net/http"
	"net/http/httptest"
	"api/v1/invoice"
	"services/route"
	"services/database"
	"database/sql/driver"
	"github.com/erikstmartin/go-testdb"
)

func beforeTest() {
	conn, err := database.OpenConnection("testdb", "")
	if err != nil {
		log.Fatal("[DB err ]: %s", err)
	}
	database.SetDb(conn)

}

func returnTwoRecords() {
	testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
		columns := []string{
			"id",
			"reference_month",
			"reference_year",
			"document",
			"description",
			"amount",
			"is_active",
			"created_at",
			"deactive_at",
		}
		result := `
		1,10,2016,069,Teste 1,12.32,1,2016-12-13 12:33:42,2016-12-14 18:59:31
		11,12,2016,6969,Teste 1,12.32,1,2016-12-15 13:24:38,2016-12-14 18:59:31
		`
		return testdb.RowsFromCSVString(columns, result), nil
	})
}

func returnNoneRecord() {
	testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
		columns := []string{
			"id",
			"reference_month",
			"reference_year",
			"document",
			"description",
			"amount",
			"is_active",
			"created_at",
			"deactive_at",
		}
		result := `
		`
		return testdb.RowsFromCSVString(columns, result), nil
	})
}

func TestGetInvoiceShouldReturnStatus404IfInvalidRequest(t *testing.T) {
	beforeTest()
	returnTwoRecords()
	data := url.Values{}
	req, err := http.NewRequest("GET", "/v1/invoicessss", bytes.NewBufferString(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer dXNlcm5hbWU6cGFzc3dvcmQ")

	rr := httptest.NewRecorder()

	router := route.NewRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := "404 page not found\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetInvoiceShouldReturnStatus404IfInvoiceNotInformed(t *testing.T) {
	beforeTest()
	returnTwoRecords()
	data := url.Values{}
	req, err := http.NewRequest("GET", "/v1/invoice/", bytes.NewBufferString(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer dXNlcm5hbWU6cGFzc3dvcmQ")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(invoice.GetInvoice)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := "Attribute invoice number is required\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetInvoiceShouldReturnStatus404IfInvoiceNotFound(t *testing.T) {
	beforeTest()
	returnNoneRecord()
	data := url.Values{}
	req, err := http.NewRequest("GET", "/v1/invoice/171", bytes.NewBufferString(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer dXNlcm5hbWU6cGFzc3dvcmQ=")

	rr := httptest.NewRecorder()

	router := route.NewRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := "Invoice [171] not found\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetInvoiceShouldReturnStatus200AndReturnValues(t *testing.T) {
	beforeTest()
	returnTwoRecords()
	data := url.Values{}
	req, err := http.NewRequest("GET", "/v1/invoice/6969", bytes.NewBufferString(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer dXNlcm5hbWU6cGFzc3dvcmQ=")

	rr := httptest.NewRecorder()

	router := route.NewRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"ID":11,"ReferenceMonth":12,"ReferenceYear":2016,"Document":"6969","Description":"Teste 1","Amount":12.32,"IsActive":true,"CreatedAt":"0001-01-01T00:00:00Z","DeactiveAt":"0001-01-01T00:00:00Z"}`+ "\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

