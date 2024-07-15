package models

import "time"

type Product struct {
	ID           uint `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time
	ProductName  string `json:"product_name"`
	SerialNumber string `json:"serial_number"`

}