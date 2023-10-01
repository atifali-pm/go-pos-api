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
	app.Patch("/cashiers/:cashier_id", controllers.UpdateCashier)
	app.Delete("/cashiers/:cashier_id", controllers.DeleteCashier)

	app.Post("/categories", controllers.CreateCategory)
	app.Get("/categories/:category_id", controllers.GetCategoryDetails)
	app.Patch("/categories/:category_id", controllers.UpdateCategory)
	app.Delete("/categories/:category_id", controllers.DeleteCategory)
	app.Get("/categories", controllers.GetCategoriesList)

	app.Post("/products", controllers.CreateProduct)
	app.Get("/products/:product_id", controllers.GetProductDetails)

}
