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
	app.Patch("/products/:product_id", controllers.UpdateProduct)
	app.Get("/products", controllers.GetProductsList)
	app.Delete("/products/:product_id", controllers.DeleteProduct)

	app.Post("/payments", controllers.CreatePayment)
	app.Get("/payments", controllers.PaymentList)
	app.Get("/payments/:payment_id", controllers.GetPaymentDetail)
	app.Delete("/payments/:payment_id", controllers.DeletePayment)
	app.Patch("/payments/:payment_id", controllers.UpdatePayment)

	app.Post("/orders", controllers.CreateOrder)
	app.Post("/orders/subtotal", controllers.SubTotalOrder)
	app.Get("/orders/check-order/:order_id", controllers.CheckOrder)
	app.Get("/orders/:order_id", controllers.OrderDetail)
	app.Get("/orders", controllers.OrderList)
	app.Get("/orders/:order_id/download", controllers.DownloadOrder)

}
