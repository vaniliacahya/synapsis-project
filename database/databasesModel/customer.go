package databasesModel

import (
	"time"
)

type Customer struct {
	Id        string     `json:"id" gorm:"primaryKey,unique"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
