package model

import (
	"time"
	"testing"
	"model"
	"services/database"
)

func beforeTest() {
	database.GetDb().AutoMigrate(&model.Invoice{})
}

func TestGetAllShoulShowAllInvoices(t *testing.T) {
	beforeTest()
	now := time.Now()
	invoice := model.Invoice{ReferenceMonth:12, ReferenceYear:2016, Document:"069", Description:"Teste 1", Amount:12.32, DeactiveAt:now}
	database.GetDb().Save(&invoice)
	row := model.GetAll()

	var invoices []model.Invoice
	database.GetDb().Find(&invoices)

	if len(row) == 0 || len(row) != len(invoices) {
		t.Errorf("Error to get all invoices")
	}

}

