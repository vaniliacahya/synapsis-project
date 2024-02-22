package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"os"
	"regexp"
	"synapsis-project/database/databasesModel"
	"synapsis-project/domain"
	"synapsis-project/helper"
	"synapsis-project/structures/request"
	"synapsis-project/structures/response"
	"time"
)

type UseCase struct {
	data domain.Data
	red  *redis.Client
}

func New(d domain.Data, r *redis.Client) *UseCase {
	return &UseCase{
		data: d,
		red:  r,
	}
}

func (u *UseCase) ListProduct(param request.ListProductRequest) (result response.LogicReturn[response.ListProduct]) {

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
	upsertCart := databasesModel.Cart{}
	addCart := body.AddCarts
	var countDataCheck int64

	//validate body
	{
		if body.IdCustomer == "" {
			result.ErrorMsg = fmt.Errorf("id_customer required")
			result.HttpErrorCode = fiber.StatusBadRequest
			return
		}

		if addCart.IdProduct == "" {
			result.ErrorMsg = fmt.Errorf("id_product required")
			result.HttpErrorCode = fiber.StatusBadRequest
			return
		}

		if addCart.Qty <= 0 {
			result.ErrorMsg = fmt.Errorf("qty required")
			result.HttpErrorCode = fiber.StatusBadRequest
			return
		}
	}

	//get price every product mentioned
	dataProduct, count, err := u.data.ListProduct(request.ListProductRequest{IdProduct: []string{addCart.IdProduct}, Limit: 1})
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	if count != 1 {
		result.ErrorMsg = fmt.Errorf("product not found")
		result.HttpErrorCode = fiber.StatusBadRequest
		return
	}

	//get cart based on Idcustomer
	dataCart, count, _, err := u.data.ListCart(body)
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	if count > 0 {
		// check if same product exist
		for _, currentCart := range dataCart {
			if currentCart.IdProduct == addCart.IdProduct {
				upsertCart = databasesModel.Cart{
					Id:           currentCart.Id,
					IdCustomer:   currentCart.IdCustomer,
					IdProduct:    currentCart.IdProduct,
					CreatedAt:    currentCart.CreatedAt,
					UpdatedAt:    &now,
					Qty:          addCart.Qty,
					PriceProduct: dataProduct[0].Price,
					TotalPrice:   dataProduct[0].Price * addCart.Qty,
				}

				//product exist in cart, then update cart
				err = u.data.UpdateCart(upsertCart)
				if err != nil {
					result.ErrorMsg = err
					result.HttpErrorCode = fiber.StatusInternalServerError
					return
				}

				break
			}
			countDataCheck++
		}
	}

	if count <= 0 || countDataCheck == count {
		//product not exist in cart, then insert
		upsertCart = databasesModel.Cart{
			Id:           uuid.New().String(),
			IdCustomer:   body.IdCustomer,
			IdProduct:    addCart.IdProduct,
			Qty:          addCart.Qty,
			PriceProduct: dataProduct[0].Price,
			TotalPrice:   dataProduct[0].Price * addCart.Qty,
		}

		err = u.data.AddCart(upsertCart)
		if err != nil {
			result.ErrorMsg = err
			result.HttpErrorCode = fiber.StatusInternalServerError
			return
		}
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

	REDISKEY := fmt.Sprintf("CART:%s", body.IdCustomer)

	//check if data exist in redis
	exist, err := u.red.Exists(context.Background(), REDISKEY).Result()
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	//if data found, then delete
	if exist == 1 {
		err = u.red.Del(context.Background(), REDISKEY).Err()
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	//insert to redis
	dataMarshall, err := json.Marshal(result.Response)
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	err = u.red.Set(context.Background(), REDISKEY, dataMarshall, 0).Err()
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	return
}

func (u *UseCase) ListCart(param request.AddCartRequest) (result response.LogicReturn[response.ListCart]) {
	//validate body
	if param.IdCustomer == "" {
		result.ErrorMsg = fmt.Errorf("id_customer required")
		result.HttpErrorCode = fiber.StatusBadRequest
		return
	}

	REDISKEY := fmt.Sprintf("CART:%s", param.IdCustomer)

	//cek if data exist in redis
	dataCartRedis, err := u.red.Get(context.Background(), REDISKEY).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			result.ErrorMsg = err
			result.HttpErrorCode = fiber.StatusInternalServerError
			return
		}

		// if redis.nil, then get data from db
		dataCart, count, total, err := u.data.ListCart(param)
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
			Products: dataCart,
			Total:    total,
		}

		//insert to redis
		dataMarshall, err := json.Marshal(result.Response)
		if err != nil {
			result.ErrorMsg = err
			result.HttpErrorCode = fiber.StatusInternalServerError
			return
		}

		err = u.red.Set(context.Background(), REDISKEY, dataMarshall, 0).Err()
		if err != nil {
			result.ErrorMsg = err
			result.HttpErrorCode = fiber.StatusInternalServerError
			return
		}

		return
	}

	err = json.Unmarshal([]byte(dataCartRedis), &result.Response)
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	return
}

func (u *UseCase) DeleteCart(param request.DeleteCartRequest) (result response.LogicReturn[response.ListCart]) {
	//validate body
	if param.IdCustomer == "" || param.Id == "" {
		result.ErrorMsg = fmt.Errorf("id and id_customer required")
		result.HttpErrorCode = fiber.StatusBadRequest
		return
	}

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

	REDISKEY := fmt.Sprintf("CART:%s", param.IdCustomer)

	//delete redis key
	err = u.red.Del(context.Background(), REDISKEY).Err()
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	//insert to redis
	dataMarshall, err := json.Marshal(result.Response)
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	err = u.red.Set(context.Background(), REDISKEY, dataMarshall, 0).Err()
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	return
}

func (u *UseCase) Order(param request.OrderRequest) (result response.LogicReturn[response.SummaryOrder]) {

	//validate body
	if param.IdCustomer == "" {
		result.ErrorMsg = fmt.Errorf("id_customer required")
		result.HttpErrorCode = fiber.StatusBadRequest
		return
	}

	cartData := u.ListCart(request.AddCartRequest{IdCustomer: param.IdCustomer})
	if cartData.ErrorMsg != nil {
		result.ErrorMsg = cartData.ErrorMsg
		result.HttpErrorCode = cartData.HttpErrorCode
		return
	}

	if cartData.Response.Count == 0 {
		result.ErrorMsg = fmt.Errorf("no cart found")
		result.HttpErrorCode = fiber.StatusBadRequest
		return
	}

	now := time.Now()
	expiredAt := now.Add(time.Hour * 24)
	adminFee := generateAdminFee()
	total := cartData.Response.Total + adminFee

	orderData := databasesModel.Order{
		Id:         uuid.New().String(),
		IdOrder:    generateIdOrder(),
		IdCustomer: param.IdCustomer,
		Subtotal:   cartData.Response.Total,
		AdminFee:   adminFee,
		Total:      total,
		ExpiredAt:  &expiredAt,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}

	err := u.data.InsertOrder(orderData)
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	REDISKEY := fmt.Sprintf("CART:%s", param.IdCustomer)

	//delete redis key
	err = u.red.Del(context.Background(), REDISKEY).Err()
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	dateExpired := expiredAt.Format("02 January 2006 at 15:04")

	result.Response = response.SummaryOrder{
		Order: orderData,
		Description: fmt.Sprintf("Your invoice %s has been successfully processed and now awaiting for payment. ", orderData.IdOrder) +
			fmt.Sprintf("Please make the payment of IDR %0.0f to the following bank account before %s to ensure your order is secured.", orderData.Total, dateExpired),
		Carts: cartData.Response.Products,
	}

	return
}

func (u *UseCase) AddCustomer(body databasesModel.Customer) (result response.LogicReturn[databasesModel.Customer]) {

	//validate request
	{
		if body.Name == "" {
			result.ErrorMsg = fmt.Errorf("name required")
			result.HttpErrorCode = fiber.StatusBadRequest
			return
		}

		if body.Username == "" {
			result.ErrorMsg = fmt.Errorf("username required")
			result.HttpErrorCode = fiber.StatusBadRequest
			return
		}

		if body.Password == "" {
			result.ErrorMsg = fmt.Errorf("password required")
			result.HttpErrorCode = fiber.StatusBadRequest
			return
		}

		validatePass := validatePassword(body.Password)
		if validatePass != nil {
			result.ErrorMsg = fmt.Errorf("invalid password format : %v", validatePass.Error())
			result.HttpErrorCode = fiber.StatusBadRequest
			return
		}
	}

	hash, err := encryptPass(body.Password)
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	now := time.Now()
	body.Id = uuid.New().String()
	body.CreatedAt = &now
	body.UpdatedAt = &now
	body.Password = string(hash)

	err = u.data.AddCustomer(body)
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	result.Response = body

	return
}

func (u *UseCase) Login(body databasesModel.Customer) (result response.LogicReturn[response.LoginResponse]) {

	//validate request
	{
		if body.Username == "" {
			result.ErrorMsg = fmt.Errorf("username required")
			result.HttpErrorCode = fiber.StatusBadRequest
			return
		}

		if body.Password == "" {
			result.ErrorMsg = fmt.Errorf("password required")
			result.HttpErrorCode = fiber.StatusBadRequest
			return
		}
	}

	dataCustomer, err := u.data.GetCustomer(body)
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	if dataCustomer.Id == "" {
		result.ErrorMsg = fmt.Errorf("username not found")
		result.HttpErrorCode = fiber.StatusNotFound
		return
	}

	//compare pass
	{
		passKey := body.Password + os.Getenv("PASS_KEY")
		err = bcrypt.CompareHashAndPassword([]byte(dataCustomer.Password), []byte(passKey))
		if err != nil {
			result.ErrorMsg = fmt.Errorf("invalid password")
			result.HttpErrorCode = fiber.StatusBadRequest
			return
		}
	}

	//generate JWT
	token, err := helper.GenerateJWT(dataCustomer.Id)
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	result.Response.Token = token

	return
}

func generateIdOrder() string {
	now := time.Now()
	rand.Seed(now.UnixNano())
	randInt := rand.Intn(100)
	dateId := now.Format("02012006")

	return fmt.Sprintf("#order%s%d", dateId, randInt)
}

func generateAdminFee() float64 {
	rand.Seed(time.Now().UnixNano())
	adminFee := rand.Intn(999-100) + 100
	return float64(adminFee)
}

func validatePassword(pass string) (err error) {
	if len(pass) < 8 {
		return fmt.Errorf("password should at least 8 character")
	}

	//password must contain uppercase
	uppercaseRegex := regexp.MustCompile("[A-Z]")
	if !uppercaseRegex.MatchString(pass) {
		return fmt.Errorf("password should contain uppercase")
	}

	//password must contain digit
	digitRegex := regexp.MustCompile("\\d")
	if !digitRegex.MatchString(pass) {
		return fmt.Errorf("password should contain digit")
	}

	return
}

func encryptPass(pass string) (hash []byte, err error) {
	godotenv.Load()
	var (
		passKey  = pass + os.Getenv("PASS_KEY")
		passByte = []byte(passKey)
		cost     = bcrypt.DefaultCost
	)

	hash, err = bcrypt.GenerateFromPassword(passByte, cost)
	if err != nil {
		return nil, err
	}

	return

}
