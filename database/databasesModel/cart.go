package databasesModel

import (
	"time"
)

type Cart struct {
	Id           string     `json:"id" gorm:"primaryKey,unique"`
	IdCustomer   string     `json:"id_customer"`
	IdProduct    string     `json:"id_product"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
	Qty          float64    `json:"qty"`
	PriceProduct float64    `json:"price"`
	TotalPrice   float64    `json:"total_price"`
}
