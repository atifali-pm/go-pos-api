package routes

import (
	"github.com/atifali-pm/go-pos-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	app.Post("/cashiers/:cashier_id/login", controllers.Login)

	app.Get("/cashiers", controllers.CashierList)
	app.Get("/cashiers/:cashier_id", controllers.GetCashierDetails)
	app.Post("/cashiers", controllers.CreateCashier)

	app.Post("/categories", controllers.CreateCategory)
	app.Get("/categories/:category_id", controllers.GetCategoryDetails)

}
