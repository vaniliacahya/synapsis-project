package usecase

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"synapsis-project/database/databasesModel"
	"synapsis-project/domain"
	"synapsis-project/structures/request"
	"synapsis-project/structures/response"
	"time"
)

type UseCase struct {
	data domain.Data
}

func New(d domain.Data) *UseCase {
	return &UseCase{
		data: d,
	}
}

func (u *UseCase) ListProduct(param request.ListProductRequest) (result response.LogicReturn[response.ListProduct]) {
	if param.Limit == 0 {
		param.Limit = 10 //default
	}

	dataProduct, count, err := u.data.ListProduct(param)
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	if count == 0 {
		result.Response.Products = []databasesModel.Product{}
		return
	}

	result.Response = response.ListProduct{
		Count:    count,
		Products: dataProduct,
	}

	return
}

func (u *UseCase) AddCart(body request.AddCartRequest) (result response.LogicReturn[response.ListCart]) {

	now := time.Now()
	idProduct := []string{}
	newInsertCart := []databasesModel.Cart{}
	newUpdateCart := []databasesModel.Cart{}
	cartMap := make(map[string]databasesModel.Cart)
	productMap := make(map[string]databasesModel.Product)
	cartExistMap := make(map[string]databasesModel.Cart)

	//validate body
	{
		if body.IdCustomer == "" {
			result.ErrorMsg = fmt.Errorf("id_customer required")
			result.HttpErrorCode = fiber.StatusBadRequest
			return
		}

		if len(body.AddCarts) <= 0 {
			result.ErrorMsg = fmt.Errorf("add_carts required")
			result.HttpErrorCode = fiber.StatusBadRequest
			return
		}

		for _, req := range body.AddCarts {
			if req.IdProduct == "" {
				result.ErrorMsg = fmt.Errorf("id_product required")
				result.HttpErrorCode = fiber.StatusBadRequest
				return
			}

			if req.Qty <= 0 {
				result.ErrorMsg = fmt.Errorf("qty required")
				result.HttpErrorCode = fiber.StatusBadRequest
				return
			}

			//needed for save idproduct used
			idProduct = append(idProduct, req.IdProduct)

			// mapping cart
			cartMap[req.IdProduct] = databasesModel.Cart{
				Id:         uuid.New().String(),
				IdCustomer: body.IdCustomer,
				IdProduct:  req.IdProduct,
				CreatedAt:  &now,
				UpdatedAt:  &now,
				Qty:        req.Qty,
			}
		}
	}

	//get price every product mentioned
	dataProduct, count, err := u.data.ListProduct(request.ListProductRequest{IdProduct: idProduct, Limit: len(idProduct)})
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	if count != int64(len(body.AddCarts)) {
		result.ErrorMsg = fmt.Errorf("product not found")
		result.HttpErrorCode = fiber.StatusBadRequest
		return
	}

	// mapping data product
	for _, p := range dataProduct {
		productMap[p.Id] = p
	}

	//get cart based on Idcustomer
	dataCart, count, _, err := u.data.ListCart(body)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	if count > 0 {
		//mapping existing cart
		for _, cart := range dataCart {
			cartExistMap[cart.IdProduct] = cart
		}
	}

	// check if same product exist
	for key, cart := range cartMap {
		//if same product exist, then update qty and price
		if _, ok := cartExistMap[key]; ok {
			newUpdateCart = append(newUpdateCart, databasesModel.Cart{
				Id:           cartExistMap[key].Id,
				IdCustomer:   cartExistMap[key].IdCustomer,
				IdProduct:    cartExistMap[key].IdProduct,
				CreatedAt:    cartExistMap[key].CreatedAt,
				UpdatedAt:    &now,
				Qty:          cartMap[key].Qty,
				PriceProduct: productMap[key].Price,
				TotalPrice:   productMap[key].Price * cartMap[key].Qty,
			})
		} else {
			//if product doesn't exist, then update price
			cart.PriceProduct = productMap[key].Price
			cart.TotalPrice = productMap[key].Price * cartMap[key].Qty
			newInsertCart = append(newInsertCart, cart)
		}
	}

	//upsert cart
	err = u.data.UpsertCart(newInsertCart, newUpdateCart)
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	// get cart based on Idcustomer
	dataCart, count, total, err := u.data.ListCart(body)
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	result.Response = response.ListCart{
		Count:    count,
		Products: dataCart,
		Total:    total,
	}

	return
}

func (u *UseCase) ListCart(param request.AddCartRequest) (result response.LogicReturn[response.ListCart]) {

	dataProduct, count, total, err := u.data.ListCart(param)
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	if count == 0 {
		result.Response.Products = []databasesModel.Cart{}
		return
	}

	result.Response = response.ListCart{
		Count:    count,
		Products: dataProduct,
		Total:    total,
	}

	return
}

func (u *UseCase) DeleteCart(param request.DeleteCartRequest) (result response.LogicReturn[response.ListCart]) {

	err := u.data.DeleteCart(param)
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	dataProduct, count, total, err := u.data.ListCart(request.AddCartRequest{IdCustomer: param.IdCustomer})
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	if count == 0 {
		result.Response.Products = []databasesModel.Cart{}
		return
	}

	result.Response = response.ListCart{
		Count:    count,
		Products: dataProduct,
		Total:    total,
	}

	return
}
