package model

import (
	"errors"
	"time"
	"services/database"
)

type Invoice struct {
	ID             uint      `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	ReferenceMonth int
	ReferenceYear  int
	Document       string    `sql:"size:14"`
	Description    string    `sql:"size:255"`
	Amount         float32   `sql:"type:decimal(10,2)"`
	IsActive       bool      `sql:"default:1"`
	CreatedAt      time.Time `sql:"type:datetime;default:current_timestamp"`
	DeactiveAt     time.Time `sql:"type:datetime"`
}

func GetAll() []Invoice {
	var invoices []Invoice
	database.GetDb().Find(&invoices)
	return invoices
}

func Delete(invoiceNumber string) error {
	var invoice Invoice
	var err error
	invoice, err = Get(invoiceNumber)
	if (err != nil) {
		return err
	}
	invoice.IsActive = false
	invoice.DeactiveAt = time.Now()
	return Save(&invoice)
}

func Update(attributes *Invoice) error {
	invoice, err := Get(attributes.Document)
	if (err != nil) {
		return err
	}
	attributes.ID = invoice.ID
	return Save(attributes)
}

func Create(attributes *Invoice) error {
	_, err := Get(attributes.Document)
	if (attributes.Document == "") {
		err = errors.New("Attribute Document containing invoice number is required")
		return err
	}
	if (err == nil) {
		err = errors.New("Invoice [" + attributes.Document + "] already exist")
		return err
	}
	return Save(attributes)
}

func Save(i *Invoice) error {
	return database.GetDb().Save(&i).Error
}

func Get(invoiceNumber string) (Invoice, error) {
	var invoice Invoice
	var err error
	if database.GetDb().Where("document = ?", invoiceNumber).Find(&invoice).RecordNotFound() {
		err = errors.New("Invoice [" + invoiceNumber + "] not found")
		return invoice, err
	}
	return invoice, nil
}

