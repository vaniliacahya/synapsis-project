package usecase

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"os"
	"regexp"
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
	if err != nil {
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
	//validate body
	if param.IdCustomer == "" {
		result.ErrorMsg = fmt.Errorf("id_customer required")
		result.HttpErrorCode = fiber.StatusBadRequest
		return
	}

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

	dateExpired := expiredAt.Format("02 January 2006 at 15:04")

	result.Response = response.SummaryOrder{
		Order: orderData,
		Description: fmt.Sprintf("Your invoice %s has been successfully processed and now awaiting for payment. ", orderData.IdOrder) +
			fmt.Sprintf("Please make the payment of IDR %0.0f to the following bank account before %s to ensure your order is secured.", orderData.Total, dateExpired),
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
	body.Password = hash

	err = u.data.AddCustomer(body)
	if err != nil {
		result.ErrorMsg = err
		result.HttpErrorCode = fiber.StatusInternalServerError
		return
	}

	result.Response = body

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

func encryptPass(pass string) (hash string, err error) {
	godotenv.Load()
	var (
		passKey  = pass + os.Getenv("PASS_KEY")
		passByte = []byte(passKey)
		cost     = bcrypt.DefaultCost
	)

	hashByte, err := bcrypt.GenerateFromPassword(passByte, cost)
	if err != nil {
		return "", err
	}

	hash = string(hashByte)
	return

}
