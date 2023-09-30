package controllers

import (
	"log"
	"strconv"
	"time"

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
		Name:      data["name"],
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
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

func UpdateCategory(c *fiber.Ctx) error {
	categoryId := c.Params("category_id")

	// Token authenticate
	if err := middleware.AuthorizeToken(c); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	// Token authenticate

	var category models.Category
	db.DB.Find(&category, "id=?", categoryId)

	if category.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Category not found against this Id",
		})
	}

	var updateCashierData models.Category
	c.BodyParser(&updateCashierData)

	if updateCashierData.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Category name is required",
			"error":   map[string]interface{}{},
		})
	}

	category.Name = updateCashierData.Name
	category.UpdatedAt = time.Now().UTC()
	db.DB.Save(&category)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    category,
	})

}

func DeleteCategory(c *fiber.Ctx) error {
	categoryId := c.Params("category_id")

	// Token authenticate
	if err := middleware.AuthorizeToken(c); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	// Token authenticate

	var category models.Category
	db.DB.Where("id=?", categoryId).First(&category)

	if category.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Category not found!",
			"error":   map[string]interface{}{},
		})
	}

	db.DB.Where("id = ?", categoryId).Delete(&category)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
	})
}

func GetCategoriesList(c *fiber.Ctx) error {
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
	var category []models.Category
	db.DB.Select("id, name").Limit(limit).Offset(skip).Find(&category).Count(&count)

	metaMap := map[string]interface{}{
		"total": count,
		"limit": limit,
		"skip":  skip,
	}

	categoriesData := map[string]interface{}{
		"categories": category,
		"meta":       metaMap,
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    categoriesData,
	})

}
