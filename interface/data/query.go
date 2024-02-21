package data

import (
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"synapsis-project/config"
	"synapsis-project/database/databasesModel"
	"synapsis-project/structures/request"
)

type Data struct {
	dbMysql gorm.DB
}

func New(dbMysql gorm.DB) *Data {
	return &Data{
		dbMysql: dbMysql,
	}
}

func (d *Data) ListProduct(param request.ListProductRequest) (res []databasesModel.Product, count int64, err error) {
	db := d.dbMysql.Where("deleted_at IS NULL")

	if len(param.IdProductCategory) > 0 {
		db = db.Where("id_product_category IN (?)", param.IdProductCategory)
	}

	if len(param.IdProduct) > 0 {
		db = db.Where("id IN (?)", param.IdProduct)
	}

	//count
	if err = db.Table(config.TableProduct).Count(&count).Error; err != nil {
		err = fmt.Errorf("get count: %v", err.Error())
		return
	}

	if param.Limit == 0 {
		param.Limit = int(count)
	}

	//select
	err = db.
		Limit(param.Limit).
		Offset(param.Offset).
		Order("name").
		Find(&res).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("get list product: %v", err.Error())
			return
		}
		err = nil
	}

	return
}

func (d *Data) ListCart(param request.AddCartRequest) (result []databasesModel.Cart, count int64, total float64, err error) {
	type customStruct struct {
		Total float64 `json:"total"`
		Count int64   `json:"count"`
	}
	data := customStruct{}

	db := d.dbMysql.Table(fmt.Sprintf("%s c", config.TableCart))

	if param.IdCustomer != "" {
		db.Where("c.id_customer", param.IdCustomer)
	}

	//count
	if err = db.Select("sum(c.total_price) as total, count(c.id) as count").Take(&data).Error; err != nil {
		err = fmt.Errorf("get total and count: %v", err.Error())
		return
	}

	total = data.Total
	count = data.Count

	//select
	err = db.Select("c.*").Limit(int(count)).
		Order("c.created_at DESC").
		Find(&result).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("get list cart: %v", err.Error())
			return
		}
		err = nil
	}

	return
}

func (d *Data) AddCart(requestInsert databasesModel.Cart) (err error) {

	if err = d.dbMysql.Create(&requestInsert).Error; err != nil {
		err = fmt.Errorf("add cart: %v", err.Error())
		return err
	}

	return
}

func (d *Data) UpdateCart(cart databasesModel.Cart) (err error) {

	if err = d.dbMysql.
		Table(fmt.Sprintf("%s c", config.TableCart)).
		Where("c.id = ?", cart.Id).
		Updates(&cart).Error; err != nil {
		err = fmt.Errorf("update cart: %v", err.Error())
		return err
	}

	return
}

func (d *Data) DeleteCart(deleteRequest request.DeleteCartRequest) (err error) {

	if err = d.dbMysql.Where("id = ?", deleteRequest.Id).
		Delete(&databasesModel.Cart{}).Error; err != nil {
		err = fmt.Errorf("delete cart: %v", err.Error())
		return
	}

	return
}

func (d *Data) InsertOrder(body databasesModel.Order) (err error) {

	return d.dbMysql.Transaction(func(tx *gorm.DB) error {

		//insert order
		if err = d.dbMysql.Create(&body).Error; err != nil {
			err = fmt.Errorf("insert order: %v", err.Error())
			return err
		}

		//delete all cart related idcustomer
		if err = d.dbMysql.Where("id_customer = ?", body.IdCustomer).Delete(&databasesModel.Cart{}).Error; err != nil {
			err = fmt.Errorf("delete all cart customer: %v", err.Error())
			return err
		}

		return nil
	})
}

func (d *Data) AddCustomer(body databasesModel.Customer) (err error) {
	//insert customer
	err = d.dbMysql.Create(&body).Error

	if err != nil {
		var sqlErr *mysql.MySQLError
		if errors.As(err, &sqlErr) {
			if sqlErr.Number == 1062 {
				err = fmt.Errorf("username is taken")
				return
			}
		}
		err = fmt.Errorf("insert customer: %v", err.Error())
		return err
	}

	return nil
}

func (d *Data) GetCustomer(body databasesModel.Customer) (res databasesModel.Customer, err error) {

	if err = d.dbMysql.Table(config.TableCustomer).
		Where("username = ?", body.Username).
		First(&res).Error; err != nil {
		err = fmt.Errorf("get customer : %v", err.Error())
		return
	}

	return
}
