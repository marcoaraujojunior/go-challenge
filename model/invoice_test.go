package model

import (
	"time"
	"testing"
	"model"
	"services/database"
)

func beforeTest() {
	database.Connect()
	database.Db.AutoMigrate(&model.Invoice{})
}

func TestGetAllShoulShowAllInvoices(t *testing.T) {
	beforeTest()
	now := time.Now()
	invoice := model.Invoice{ReferenceMonth:12, ReferenceYear:2016, Document:"069", Description:"Teste 1", Amount:12.32, DeactiveAt:now}
	database.Db.Save(&invoice)
	row := model.GetAll().Row()

	var test model.Invoice
	row.Scan(&test)
	t.Errorf(test.Document)
}

