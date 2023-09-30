package controllers

import (
	"log"

	db "github.com/atifali-pm/go-pos-api/config"
	"github.com/atifali-pm/go-pos-api/middleware"
	"github.com/atifali-pm/go-pos-api/models"
	"github.com/gofiber/fiber/v2"
)

func CreateCategory(c *fiber.Ctx) error {
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		log.Fatalf("category error in post requires %v", err)
	}

	if data["name"] == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Category name is required",
			"error":   map[string]interface{}{},
		})
	}

	category := models.Category{
		Name: data["name"],
	}

	db.DB.Create(&category)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    category,
	})
}

func GetCategoryDetails(c *fiber.Ctx) error {

	categoryId := c.Params("category_id")

	headerToken := c.Get("Authorization")

	if headerToken == "" {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   map[string]interface{}{},
		})
	}

	if err := middleware.AuthenticateToken(middleware.SplitToken(headerToken)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   map[string]interface{}{},
		})
	}

	var category models.Category
	db.DB.Select("id, name").Where("id=?", categoryId).First(&category)

	categoryData := make(map[string]interface{})
	categoryData["categoryId"] = category.Id
	categoryData["name"] = category.Name

	if category.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"success": true,
			"message": "No category found!",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    categoryData,
	})
}
