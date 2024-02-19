package dispatch

import (
	"github.com/gofiber/fiber/v2"
	"synapsis-project/domain"
)

func Routes(app *fiber.App, handler domain.Handler) {
	//product
	app.Get("/products", handler.ListProduct())

	//cart
	cart := app.Group("/carts")
	cart.Post("", handler.AddCart())
	cart.Get("", handler.ListCart())
	cart.Delete("/:id", handler.DeleteCart())

	//checkout order
	app.Post("/checkout", handler.Order())

	//customer
	app.Post("/register", handler.Register())
	app.Post("/login", handler.Login())
}
