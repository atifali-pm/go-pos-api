package main

import (
	"fmt"

	db "github.com/atifali-pm/go-pos-api/config"
	"github.com/atifali-pm/go-pos-api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	fmt.Println("here")

	db.Connect()

	app := fiber.New()

	app.Use(cors.New())

	routes.Setup(app)

	app.Listen(":3030")

}
