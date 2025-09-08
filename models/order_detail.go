package models

import (
	// "time"
)

type Order_Detail struct {
	Id          int64  `gorm:"primaryKey" json:"id"`
	Order_id int `json:"order_id"`
	product_id   int `json:"product_id"`
}
