package controllers

import (
	"log"
	"strconv"

	db "github.com/atifali-pm/go-pos-api/config"
	"github.com/atifali-pm/go-pos-api/middleware"
	"github.com/atifali-pm/go-pos-api/models"

	"github.com/gofiber/fiber/v2"
)

func CreateCashier(c *fiber.Ctx) error {
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		log.Fatalf("Cashier not registered, fatal error %v", err)
	}

	if data["name"] == "" || data["passcode"] == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Cashier name is required",
			"error":   map[string]interface{}{},
		})
	}

	cashier := models.Cashier{
		Name:     data["name"],
		Passcode: data["passcode"],
	}
	db.DB.Create(&cashier)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    cashier,
	})
}

func GetCashierDetails(c *fiber.Ctx) error {
	cashierId := c.Params("cashier_id")

	var cashier models.Cashier
	db.DB.Select("id, name").Where("id=?", cashierId).First(&cashier)

	cashierData := make(map[string]interface{})
	cashierData["id"] = cashier.Id
	cashierData["name"] = cashier.Name

	if cashier.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier not found!",
			"error":   map[string]interface{}{},
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    cashierData,
	})
}

type Cashiers struct {
	Id   uint   `json:"cashierId"`
	Name string `json:"name"`
}

func CashierList(c *fiber.Ctx) error {

	limitParam := c.Query("limit")

	// Set default limit to 10 if not provided or invalid
	limit := 10
	if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
		limit = parsedLimit
	}

	skip, _ := strconv.Atoi(c.Query("skip"))

	var count int64
	var cashier []Cashiers

	db.DB.Select("*").Limit(limit).Offset(skip).Find(&cashier).Count(&count)
	metaMap := map[string]interface{}{
		"total": count,
		"limit": limit,
		"skip":  skip,
	}

	cashiersData := map[string]interface{}{
		"cashiers": cashier,
		"meta":     metaMap,
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    cashiersData,
	})

}

func UpdateCashier(c *fiber.Ctx) error {
	cashierId := c.Params("cashier_id")

	// Token authenticate
	if err := middleware.AuthorizeToken(c); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	// Token authenticate

	var cashier models.Cashier

	db.DB.Find(&cashier, "id = ?", cashierId)

	if cashier.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier not found!",
		})
	}

	var updateCashierData models.Cashier
	c.BodyParser(&updateCashierData)

	if updateCashierData.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Cashier name is required",
			"error":   map[string]interface{}{},
		})
	}

	if updateCashierData.Passcode == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": true,
			"message": "Cashier passcode is required",
			"error":   map[string]interface{}{},
		})
	}

	cashier.Name = updateCashierData.Name
	cashier.Passcode = updateCashierData.Passcode

	db.DB.Save(&cashier)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    cashier,
	})
}

func DeleteCashier(c *fiber.Ctx) error {
	cashierId := c.Params("cashier_id")

	if err := middleware.AuthorizeToken(c); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	var cashier models.Cashier
	db.DB.Where("id = ?", cashierId).First(&cashier)

	if cashier.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "cashier not found!",
			"error":   map[string]interface{}{},
		})
	}

	db.DB.Where("id = ?", cashierId).Delete(&cashier)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
	})
}
