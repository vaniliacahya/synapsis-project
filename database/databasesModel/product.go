package databasesModel

import (
	"time"
)

type Product struct {
	Id                string     `json:"id" gorm:"primaryKey,unique"`
	Name              string     `json:"name"`
	IdProductCategory string     `json:"id_product_category"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
	Price             float64    `json:"price"`
}
