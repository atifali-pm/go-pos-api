package routes

import (
	"github.com/atifali-pm/go-pos-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	app.Get("/cashiers", controllers.CashierList)
	app.Post("/cashiers", controllers.CreateCashier)

}
