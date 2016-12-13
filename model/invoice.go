package model

import (
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

