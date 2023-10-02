package controllers

import (
	"log"
	"strconv"

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

func PaymentList(c *fiber.Ctx) error {
	// Token authenticate
	if err := middleware.AuthorizeToken(c); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	// Token authenticate

	limitParam := c.Query("limit")

	// Set default limit to 10 if not provided or invalid
	limit := 10
	if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
		limit = parsedLimit
	}

	skip, _ := strconv.Atoi(c.Query("skip"))

	var count int64
	var payment []models.Payment
	db.DB.Select("*").Limit(limit).Offset(skip).Find(&payment).Count(&count)
	metaMap := map[string]interface{}{
		"total":  count,
		"limiit": limit,
		"skip":   skip,
	}

	paymentData := map[string]interface{}{
		"payments": payment,
		"meta":     metaMap,
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    paymentData,
	})

}

func GetPaymentDetail(c *fiber.Ctx) error {
	// Token authenticate
	if err := middleware.AuthorizeToken(c); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	// Token authenticate

	paymentId := c.Params("payment_id")
	var payment models.Payment
	db.DB.Where("id=?", paymentId).First(&payment)

	if payment.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Payment not found!",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"data":    payment,
	})
}

func DeletePayment(c *fiber.Ctx) error {
	// Token authenticate
	if err := middleware.AuthorizeToken(c); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	// Token authenticate

	paymentId := c.Params("payment_id")
	var payment models.Payment
	db.DB.First(&payment, paymentId)

	if payment.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "payment not found!",
		})
	}

	result := db.DB.Delete(&payment)
	if result.RowsAffected <= 0 {
		return c.Status(fiber.StatusNotModified).JSON(fiber.Map{
			"success": false,
			"message": "Payment removing failed",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Deleted successfully!",
	})

}

func UpdatePayment(c *fiber.Ctx) error {
	// Token authenticate
	if err := middleware.AuthorizeToken(c); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	// Token authenticate
	paymentId := c.Params("payment_id")

	var payment models.Payment
	db.DB.Find(&payment, "id = ?", paymentId)

	if payment.Id <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Payment not found",
			"error":   map[string]interface{}{},
		})
	}

	var updatePayment models.Payment
	c.BodyParser(&updatePayment)

	if updatePayment.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Payment name is required",
			"error":   map[string]interface{}{},
		})
	}

	if updatePayment.Type == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Payment type is required",
			"error":   map[string]interface{}{},
		})
	}

	var paymentTypeId int
	if updatePayment.Type == "CASH" {
		paymentTypeId = 1
	}
	if updatePayment.Type == "E-WALLET" {
		paymentTypeId = 2
	}
	if updatePayment.Type == "EDC" {
		paymentTypeId = 3
	}

	payment.Name = updatePayment.Name
	payment.Type = updatePayment.Type
	payment.PaymentTypeId = paymentTypeId
	payment.Logo = updatePayment.Logo

	db.DB.Save(&payment)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Payment updated!",
		"data":    payment,
	})

}
