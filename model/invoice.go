package model

import (
	"time"
)

type Invoice struct {
	ID             uint    `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	ReferenceMonth int
	ReferenceYear  int
	Document       string  `sql:"size:14"`
	Description    string  `sql:"size:255"`
	Amount         float32 `sql:"type:decimal(10,2)"`
	IsActive       bool
	CreatedAt      time.Time
	DeactiveAt     time.Time
}
