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

	if data.CategoryId == 0 || data.Name == "" || data.Image == "" || data.Stock <= 0 || data.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": true,
			"message": "Fields are required",
		})
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
		Id       int      `json:"id"`
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

func UpdateProduct(c *fiber.Ctx) error {

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

	db.DB.Find(&product, "id = ?", productId)

	if product.Name == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": true,
			"message": "Product not found",
		})
	}

	// fmt.Println("--------------------------------------->")
	// fmt.Println("------------Product ----------->", product)
	// fmt.Println("--------------------------------------->")

	var updateProductData models.Product
	c.BodyParser(&updateProductData)

	if updateProductData.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Product name is required",
			"error":   map[string]interface{}{},
		})
	}

	if updateProductData.Name != "" {
		product.Name = updateProductData.Name
	}
	if updateProductData.Price > 0 {
		product.Price = updateProductData.Price
	}
	if updateProductData.CategoryId > 0 {
		product.CategoryId = updateProductData.CategoryId
	}
	if updateProductData.Image != "" {
		product.Image = updateProductData.Image
	}
	if updateProductData.Stock > 0 {
		product.Stock = updateProductData.Stock
	}

	db.DB.Save(&product)

	return c.Status(200).JSON(fiber.Map{
		"success":  true,
		"messsage": "success",
		"data":     product,
	})

}

func GetProductsList(c *fiber.Ctx) error {
	// Token authenticate
	if err := middleware.AuthorizeToken(c); err != nil {
		return c.Status(403).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"status":  fiber.StatusForbidden,
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

	var products []models.Product

	productsRes := make([]*models.ProductResult, 0)

	var count int64

	db.DB.Limit(limit).Offset(skip).Find(&products).Count(&count)

	var category models.CategoryResult
	var discountResult models.DiscountResult

	for i := 0; i < len(products); i++ {
		db.DB.Table("categories").Select("id, name").Where("id = ?", products[i].CategoryId).Find(&category)
		db.DB.Table("discounts").Select("*").Where("id = ?", products[i].DiscountId).Limit(limit).Offset(skip).Find(&discountResult).Count(&count)

		count = int64(len(products))

		productsRes = append(productsRes, &models.ProductResult{
			Id:       products[i].Id,
			Sku:      products[i].Sku,
			Name:     products[i].Name,
			Stock:    products[i].Stock,
			Price:    products[i].Price,
			Image:    products[i].Image,
			Category: category,
			Discount: discountResult,
		})

	}

	meta := map[string]interface{}{
		"total": count,
		"limit": limit,
		"skip":  skip,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    productsRes,
		"meta":    meta,
	})

}

func DeleteProduct(c *fiber.Ctx) error {
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

	db.DB.First(&product, productId)
	if product.Id <= 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"messgae": "product not found!",
		})
	}

	db.DB.Delete(&product)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "product is removed!",
	})
}
