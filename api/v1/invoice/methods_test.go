package invoice_test

import(
	"log"
	"bytes"
	"strings"
	"testing"
	"net/url"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"api/v1/invoice"
	"services/route"
	"services/database"
	"database/sql/driver"
	"github.com/erikstmartin/go-testdb"
)

type testResult struct {
	lastId       int64
	affectedRows int64
}

func (r testResult) LastInsertId() (int64, error) {
	return r.lastId, nil
}

func (r testResult) RowsAffected() (int64, error) {
	return r.affectedRows, nil
}

func beforeTest() {
	conn, err := database.OpenConnection("testdb", "")
	if err != nil {
		log.Fatal("[DB err ]: %s", err)
	}
	testdb.EnableTimeParsing(true)
	database.SetDb(conn)
}

func setQueryFunc(result string, columns []string) {
	testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
		return testdb.RowsFromCSVString(columns, result), nil
	})
}

func selectAllColumns() []string {
	return []string{
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
}

func returnInvalidRecord() {
	result := `
	4,7,2016,69,Teste 69,12.31,1,2016-12-20 11:15:37,0001-01-01T00:00:00Z
	`
	setQueryFunc(result, selectAllColumns())
}

func returnOneRecord() {
	result := `
	4,7,2016,69,Teste 69,12.31,1,2016-12-20T11:15:37Z,0001-01-01T00:00:00Z
	`
	setQueryFunc(result, selectAllColumns())
}

func returnTwoRecords() {
	result := `
	1,10,2016,069,Teste 1,12.32,1,2016-12-20T11:15:37Z,0001-01-01T00:00:00Z
	11,12,2016,6969,Teste 1,12.32,1,2016-12-20T11:15:37Z,0001-01-01T00:00:00Z
	`
	setQueryFunc(result, selectAllColumns())
}

func countFiveInvoices() {
	columns := []string{"count"}
	result := `5`
	setQueryFunc(result, columns)
}

func returnFiveRecords() {
	result := `
	1,10,2016,069,Teste 1,12.32,1,2016-12-20T11:15:37Z,0001-01-01T00:00:00Z
	2,11,2016,696,Teste 1,12.32,1,2016-12-20T11:15:37Z,0001-01-01T00:00:00Z
	3,8,2016,6,Teste 1,12.32,1,2016-12-20T11:15:37Z,0001-01-01T00:00:00Z
	10,9,2016,969,Teste 1,12.32,1,2016-12-20T11:15:37Z,0001-01-01T00:00:00Z
	11,12,2016,6969,Teste 1,12.32,1,2016-12-20T11:15:37Z,0001-01-01T00:00:00Z
	`
	setQueryFunc(result, selectAllColumns())
}

func returnNoneRecord() {
	setQueryFunc(``, selectAllColumns())
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

	expected := `{"ID":11,"ReferenceMonth":12,"ReferenceYear":2016,"Document":"6969","Description":"Teste 1","Amount":12.32,"IsActive":true,"CreatedAt":"2016-12-20T11:15:37Z","DeactiveAt":"0001-01-01T00:00:00Z"}`+ "\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestListInvoicesShouldReturnStatus404IfInvalidRequest(t *testing.T) {
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

func TestListInvoicesShouldReturnStatus200AndReturnValues(t *testing.T) {
	beforeTest()
	returnTwoRecords()
	data := url.Values{}
	req, err := http.NewRequest("GET", "/v1/invoices", bytes.NewBufferString(data.Encode()))
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

	expected := `[{"ID":1,"ReferenceMonth":10,"ReferenceYear":2016,"Document":"069","Description":"Teste 1","Amount":12.32,"IsActive":true,"CreatedAt":"2016-12-20T11:15:37Z","DeactiveAt":"0001-01-01T00:00:00Z"},{"ID":11,"ReferenceMonth":12,"ReferenceYear":2016,"Document":"6969","Description":"Teste 1","Amount":12.32,"IsActive":true,"CreatedAt":"2016-12-20T11:15:37Z","DeactiveAt":"0001-01-01T00:00:00Z"}]`+ "\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestListInvoicesShouldReturnStatus200AndDoNotReturnValues(t *testing.T) {
	beforeTest()
	returnNoneRecord()
	data := url.Values{}
	req, err := http.NewRequest("GET", "/v1/invoices", bytes.NewBufferString(data.Encode()))
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

	expected := `[]`+"\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestListInvoicesShouldReturnStatus200AndFourLinksHeader(t *testing.T) {
	beforeTest()
	returnFiveRecords()
	countFiveInvoices()

	data := url.Values{}
	req, err := http.NewRequest("GET", "/v1/invoices?page=3&per_page=1", bytes.NewBufferString(data.Encode()))
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

	links := rr.Header()["Link"][0]

	expectedLinkNext := `</v1/invoices?page=4&per_page=1> ; rel="next"`
	if !strings.ContainsAny(links, expectedLinkNext) {
		t.Errorf("handler returned unexpected link header: got %v want %v", links, expectedLinkNext)
	}

	expectedLinkLast := `</v1/invoices?page=5&per_page=1> ; rel="last"`
	if !strings.ContainsAny(links, expectedLinkLast) {
		t.Errorf("handler returned unexpected link header: got %v want %v", links, expectedLinkLast)
	}

	expectedLinkFirst := `</v1/invoices?page=1&per_page=1> ; rel="first"`
	if !strings.ContainsAny(links, expectedLinkFirst) {
		t.Errorf("handler returned unexpected link header: got |%v| want |%v|", links, expectedLinkFirst)
	}

	expectedLinkPrev := `</v1/invoices?page=2&per_page=1> ; rel="prev"`
	if !strings.ContainsAny(links, expectedLinkPrev) {
		t.Errorf("handler returned unexpected link header: got [%v] want [%v]", links, expectedLinkPrev)
	}
}

func TestUpdateInvoiceShouldReturnStatus500OnError(t *testing.T) {
	beforeTest()
	returnFiveRecords()
	countFiveInvoices()

	toUpdate := map[string]interface{}{
		"referencemonth": 8,
		"ReferenceYear": 2016,
		"Document": "69",
		"Description": "Teste 69",
		"Amount": 69.69,
	}
	body, _ := json.Marshal(toUpdate)

	req, err := http.NewRequest("PUT", "/v1/invoice/69", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer dXNlcm5hbWU6cGFzc3dvcmQ=")

	rr := httptest.NewRecorder()

	router := route.NewRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

}

func TestUpdateInvoiceShouldReturnStatus204(t *testing.T) {
	beforeTest()
	returnOneRecord()

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		if args[0] == "7" {
			return testResult{1, 1}, nil
		}
		return testResult{1, 0}, nil
	})

	toUpdate := map[string]interface{}{
		"referencemonth": 7,
		"ReferenceYear": 2016,
		"Document": "69",
		"Description": "Teste 69",
		"Amount": 69.69,
	}
	body, _ := json.Marshal(toUpdate)

	req, err := http.NewRequest("PUT", "/v1/invoice/69", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer dXNlcm5hbWU6cGFzc3dvcmQ=")

	rr := httptest.NewRecorder()

	router := route.NewRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

}

func TestUpdateInvoiceShouldReturnStatus500WhenInvalidJson(t *testing.T) {
	beforeTest()
	returnOneRecord()

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		if args[0] == "7" {
			return testResult{1, 1}, nil
		}
		return testResult{1, 0}, nil
	})

	body, _ := json.Marshal(`invalid json`)

	req, err := http.NewRequest("PUT", "/v1/invoice/69", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer dXNlcm5hbWU6cGFzc3dvcmQ=")

	rr := httptest.NewRecorder()

	router := route.NewRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}
	expected := `{"Value":"string","Type":{},"Offset":14}` + "\n" + `Invalid month`+"\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

}

func TestDeleteInvoiceShouldReturnStatus204(t *testing.T) {
	beforeTest()
	returnOneRecord()

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		if args[0] == "7" {
			return testResult{1, 1}, nil
		}
		return testResult{1, 0}, nil
	})

	data := url.Values{}
	req, err := http.NewRequest("DELETE", "/v1/invoice/69", bytes.NewBufferString(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer dXNlcm5hbWU6cGFzc3dvcmQ=")

	rr := httptest.NewRecorder()

	router := route.NewRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

}

func TestDeleteInvoiceShouldReturnStatus500(t *testing.T) {
	beforeTest()
	returnInvalidRecord()

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		if args[0] == "7" {
			return testResult{1, 1}, nil
		}
		return testResult{1, 0}, nil
	})

	data := url.Values{}
	req, err := http.NewRequest("DELETE", "/v1/invoice/69", bytes.NewBufferString(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer dXNlcm5hbWU6cGFzc3dvcmQ=")

	rr := httptest.NewRecorder()

	router := route.NewRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

}

func TestCreateInvoiceShouldReturnStatus204(t *testing.T) {
	beforeTest()
	returnNoneRecord()

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		if args[0] == "7" {
			return testResult{1, 1}, nil
		}
		return testResult{1, 0}, nil
	})

	toCreate := map[string]interface{}{
		"referencemonth": 5,
		"ReferenceYear": 2016,
		"Document": "696969",
		"Description": "Teste 6969696",
		"Amount": 69.69,
	}
	body, _ := json.Marshal(toCreate)

	req, err := http.NewRequest("POST", "/v1/invoice", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer dXNlcm5hbWU6cGFzc3dvcmQ=")

	rr := httptest.NewRecorder()

	router := route.NewRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

}

func TestCreateInvoiceShouldReturnStatus500WhenInvalidJson(t *testing.T) {
	beforeTest()
	returnNoneRecord()

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		if args[0] == "7" {
			return testResult{1, 1}, nil
		}
		return testResult{1, 0}, nil
	})

	body, _ := json.Marshal(`invalid json`)

	req, err := http.NewRequest("POST", "/v1/invoice", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer dXNlcm5hbWU6cGFzc3dvcmQ=")

	rr := httptest.NewRecorder()

	router := route.NewRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}

}

func TestCreateInvoiceShouldReturnStatus500WhenInvoiceAlreadyExist(t *testing.T) {
	beforeTest()
	returnOneRecord()

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		if args[0] == "7" {
			return testResult{1, 1}, nil
		}
		return testResult{1, 0}, nil
	})

	toCreate := map[string]interface{}{
		"referencemonth": 5,
		"ReferenceYear": 2016,
		"Document": "69",
		"Description": "Teste 69",
		"Amount": 69.69,
	}
	body, _ := json.Marshal(toCreate)

	req, err := http.NewRequest("POST", "/v1/invoice", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer dXNlcm5hbWU6cGFzc3dvcmQ=")

	rr := httptest.NewRecorder()

	router := route.NewRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

}

