package models

import (
	"time"
)

type Order struct {
	Id          int64  `gorm:"primaryKey" json:"id"`
	Customer string `gorm:"type:varchar(300)" json:"customer"`
	Date   time.Time  `gorm:"type:text" json:"date"`
	Status   string `gorm:"type:text" json:"status"`
}
