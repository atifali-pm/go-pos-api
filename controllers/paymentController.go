package controllers

import (
	"log"

	db "github.com/atifali-pm/go-pos-api/config"
	"github.com/atifali-pm/go-pos-api/middleware"
	"github.com/atifali-pm/go-pos-api/models"
	"github.com/gofiber/fiber/v2"
)

func CreatePayment(c *fiber.Ctx) error {

	// Token authenticate
	if err := middleware.AuthorizeToken(c); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	// Token authenticate

	var data map[string]string
	paymentError := c.BodyParser(&data)

	if paymentError != nil {
		log.Fatalf("Payment error %v", paymentError)

	}

	if data["name"] == "" || data["type"] == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Payment name and type is required",
			"error":   map[string]interface{}{},
		})
	}

	var paymentTypes models.PaymentType
	db.DB.Where("name", data["type"]).Find(&paymentTypes)

	payment := models.Payment{
		Name:          data["name"],
		Type:          data["type"],
		PaymentTypeId: int(paymentTypes.Id),
		Logo:          data["logo"],
	}

	db.DB.Create(&payment)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    payment,
	})

}
