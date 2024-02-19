package main

import (
	"fmt"
	"github.com/google/uuid"
	"synapsis-project/config"
	"synapsis-project/database"
	"synapsis-project/database/databasesModel"
	"time"
)

func main() {
	//load connection sql
	dbCfgMysql := config.GetDBConfig()
	dbMysql := database.InitMySqlDB(dbCfgMysql)

	now := time.Now()
	product := []databasesModel.Product{}

	//seeder product_category
	productCategory := []databasesModel.ProductCategory{
		{
			Id:        uuid.New().String(),
			Name:      "Clothes",
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		{
			Id:        uuid.New().String(),
			Name:      "Bag",
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}

	for _, category := range productCategory {
		if category.Name == "Clothes" {
			product = append(product,
				databasesModel.Product{
					Id:                uuid.New().String(),
					Name:              "White T-Shirt",
					IdProductCategory: category.Id,
					CreatedAt:         &now,
					UpdatedAt:         &now,
					Price:             200_000},
				databasesModel.Product{
					Id:                uuid.New().String(),
					Name:              "Black Dress",
					IdProductCategory: category.Id,
					CreatedAt:         &now,
					UpdatedAt:         &now,
					Price:             350_000},
			)
		} else if category.Name == "Bag" {
			product = append(product,
				databasesModel.Product{
					Id:                uuid.New().String(),
					Name:              "Mini Sling Bag",
					IdProductCategory: category.Id,
					CreatedAt:         &now,
					UpdatedAt:         &now,
					Price:             120_000},
				databasesModel.Product{
					Id:                uuid.New().String(),
					Name:              "Medium Shoulder Bag",
					IdProductCategory: category.Id,
					CreatedAt:         &now,
					UpdatedAt:         &now,
					Price:             550_000},
			)
		}
	}

	//customer := &databasesModel.Customer{
	//	Id:        uuid.New().String(),
	//	Name:      "Vanilia",
	//	CreatedAt: &now,
	//	UpdatedAt: &now,
	//}

	//bulk insert product_category
	if err := dbMysql.Create(&productCategory).Error; err != nil {
		fmt.Println("error insert product category")
	}

	//bulk insert product
	if err := dbMysql.Create(&product).Error; err != nil {
		fmt.Println("error insert product")
	}

	//insert customer
	//if err := dbMysql.Create(&customer).Error; err != nil {
	//	fmt.Println("error insert customer")
	//}

	fmt.Println("success seeding data")
}
