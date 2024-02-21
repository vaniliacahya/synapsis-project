package databasesModel

import "time"

type Order struct {
	Id         string     `json:"id" gorm:"primaryKey,unique"`
	IdOrder    string     `json:"id_order"`
	IdCustomer string     `json:"id_customer"`
	Total      float64    `json:"total"`
	Paid       bool       `json:"paid"`
	Expired    bool       `json:"expired"`
	ExpiredAt  *time.Time `json:"expired_at"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
	AdminFee   float64    `json:"admin_fee"`
	Subtotal   float64    `json:"subtotal"`
}
