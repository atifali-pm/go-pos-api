package controllers

import (
	"fmt"
	"log"
	"strconv"

	db "github.com/atifali-pm/go-pos-api/config"
	"github.com/atifali-pm/go-pos-api/middleware"
	"github.com/atifali-pm/go-pos-api/models"
	"github.com/gofiber/fiber/v2"
)

type Products struct {
	Products     models.Product
	CategoriesId string `json:"categories_Id"`
}
type ProductDiscount struct {
	Id         int      `json:"id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Sku        string   `json:"sku"`
	Name       string   `json:"name"`
	Stock      int      `json:"stock"`
	Price      int      `json:"price"`
	Image      string   `json:"image"`
	CategoryId int      `json:"categoryId"`
	Discount   Discount `json:"discount"`
}
type Discount struct {
	Qty       int    `json:"qty"`
	Type      string `json:"type"`
	Result    int    `json:"result"`
	ExpiredAt int    `json:"expiredAt"`
}

type Category struct {
	Name string `json:"name"`
}

func CreateProduct(c *fiber.Ctx) error {
	var data ProductDiscount

	// Token authenticate
	if err := middleware.AuthorizeToken(c); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	// Token authenticate

	err := c.BodyParser(&data)
	if err != nil {
		log.Fatalf("Product error in post request %v", err)
	}

	// var p []models.Product
	// db.DB.Find(&p)

	// Creating discount
	discount := models.Discount{
		Qty:       data.Discount.Qty,
		Type:      data.Discount.Type,
		Result:    data.Discount.Result,
		ExpiredAt: data.Discount.ExpiredAt,
	}
	db.DB.Create(&discount)

	//Creating Product
	product := models.Product{
		Name:       data.Name,
		Image:      data.Image,
		CategoryId: data.CategoryId,
		DiscountId: discount.Id,
		Price:      data.Price,
		Stock:      data.Stock,
	}
	db.DB.Create(&product)

	sku := "SK00" + strconv.Itoa(product.Id)
	db.DB.Table("products").Where("id = ?", product.Id).Update("sku", &sku)

	fmt.Println("--------------------------------------->")
	fmt.Println("------------Product Creation Done----------->", product.Id)
	fmt.Println("--------------------------------------->")

	product.Sku = sku

	return c.Status(200).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"success": true,
		"message": "Product created!",
		"data":    product,
	})

}

func GetProductDetails(c *fiber.Ctx) error {
	// Token authenticate
	if err := middleware.AuthorizeToken(c); err != nil {
		return c.Status(403).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"status":  fiber.StatusForbidden,
		})
	}
	// Token authenticate

	productId := c.Params("product_id")

	var product models.Product
	db.DB.Where("id=?", productId).Find(&product)

	if product.Id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"status":  fiber.StatusNotFound,
			"message": "Product not found!",
			"error":   map[string]interface{}{},
		})
	}

	var category models.Category
	var discount models.Discount

	db.DB.Where("id = ?", product.CategoryId).Find(&category)
	db.DB.Where("id = ?", product.DiscountId).Find(&discount)

	type ProductDiscount struct {
		Id       int      `json:"id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
		Sku      string   `json:"sku"`
		Name     string   `json:"name"`
		Stock    int      `json:"stock"`
		Price    int      `json:"price"`
		Image    string   `json:"image"`
		Category Category `json:"Category"`
		Discount Discount `json:"Discount"`
	}

	productResponse := ProductDiscount{
		Id:    product.Id,
		Sku:   product.Sku,
		Name:  product.Name,
		Stock: product.Stock,
		Price: product.Price,
		Image: product.Image,
		Category: Category{
			Name: category.Name,
		},
		Discount: Discount{
			Qty:       discount.Qty,
			Type:      discount.Type,
			Result:    discount.Result,
			ExpiredAt: discount.ExpiredAt,
		},
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    productResponse,
	})

}
