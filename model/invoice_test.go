package model

import (
	"testing"
	"model"
	"services/database"
	"database/sql/driver"
	"github.com/erikstmartin/go-testdb"
	"log"
)

func beforeTest() {
	conn, err := database.OpenConnection("testdb", "")
	if err != nil {
		log.Fatal("[DB err ]: %s", err)
	}
	database.SetDb(conn)
}

func TestGetAllShoulShowAllInvoices(t *testing.T) {
	beforeTest()

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

	rows := model.GetAll()

	if (rows[0].Document != "069") || len(rows) != 2 || (rows[1].Document != "6969"){
		t.Errorf("Unexcepted result returned")
	}

}

