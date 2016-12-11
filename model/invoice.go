package model

import (
	"time"
	"github.com/jinzhu/gorm"
	"github.com/marcoaraujojunior/go-challenge/database"
)

type Invoice struct {
	ID             uint      `json:"-";sql:"AUTO_INCREMENT" gorm:"primary_key"`
	ReferenceMonth int       `json:"reference_month"`
	ReferenceYear  int       `json:"reference_year"`
	Document       string    `json:"document";sql:"size:14"`
	Description    string    `json:"description";sql:"size:255"`
	Amount         float32   `json:"amount";sql:"type:decimal(10,2)"`
	IsActive       bool      `json:"is_active";sql:"default:1"`
	CreatedAt      time.Time `json:"created_at";sql:"type:datetime;default:current_timestamp"`
	DeactiveAt     time.Time `json:"deactive_at,omitempty";sql:"type:datetime"`
}

func GetAll() *gorm.DB {
	invoices := []Invoice{}
	return database.Db.Find(&invoices)
}
