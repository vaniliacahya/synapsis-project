package helper

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"os"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	IdCustomer string `json:"id_customer"`
	jwt.StandardClaims
}

func GenerateJWT(idcustomer string) (tokenString string, err error) {
	claims := &Claims{
		IdCustomer: idcustomer,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	//create JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtSecret)
	if err != nil {
		err = fmt.Errorf("error generate jwt : %v", err.Error())
		return
	}

	return
}

func ExtractData(ctx *fiber.Ctx) (claims Claims, err error) {
	tokenString := ctx.Get("Authorization")

	if len(tokenString) < 7 {
		return Claims{}, fmt.Errorf("bearer token required")
	}

	tokenwobearer := tokenString[7:]

	token, err := jwt.ParseWithClaims(tokenwobearer, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, err
	})

	if err != nil {
		return
	}

	if claim, ok := token.Claims.(*Claims); ok && token.Valid {
		return *claim, nil
	}

	return Claims{}, fmt.Errorf("error extract token")
}
