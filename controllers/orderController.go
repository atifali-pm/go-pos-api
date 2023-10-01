package controllers

import (
	"math/rand"
	"strconv"
	"time"

	db "github.com/atifali-pm/go-pos-api/config"
	"github.com/atifali-pm/go-pos-api/middleware"
	"github.com/atifali-pm/go-pos-api/models"
	"github.com/gofiber/fiber/v2"
)

func CreateOrder(c *fiber.Ctx) error {
	// Token authenticate
	if err := middleware.AuthorizeToken(c); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	// Token authenticate

	type products struct {
		ProductId int `json:"productId"`
		Quantity  int `json:"qty"`
	}

	body := struct {
		PaymentTypeId int        `json:"PaymentTypeId"`
		TotalPaid     int        `json:"totalPaid"`
		Products      []products `json:"products"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Empty Body",
			"error":   map[string]interface{}{},
		})
	}

	Prodresponse := make([]*models.ProductResponseOrder, 0)

	var TotalInvoicePrice = struct {
		ttprice int
	}{}

	productsIds := ""
	quantities := ""

	for _, v := range body.Products {
		totalPrice := 0
		productsIds = productsIds + "," + strconv.Itoa(v.ProductId)
		quantities = quantities + "," + strconv.Itoa(v.Quantity)

		prods := models.ProductOrder{}
		var discount models.Discount
		db.DB.Table("products").Where("id = ?", v.ProductId).First(&prods)
		db.DB.Where("id = ?", prods.DiscountId).Find(&discount)
		discCount := 0

		if discount.Type == "Buy_N" {
			totalPrice = prods.Price * v.Quantity
			discCount = totalPrice - discount.Result
			TotalInvoicePrice.ttprice = TotalInvoicePrice.ttprice + discCount
		}

		if discount.Type == "PERCENT" {
			totalPrice = prods.Price * v.Quantity
			percentage := totalPrice * discount.Result / 100
			discCount = totalPrice - percentage
			TotalInvoicePrice.ttprice = TotalInvoicePrice.ttprice + discCount
		}

		Prodresponse = append(Prodresponse, &models.ProductResponseOrder{
			ProductId:        prods.Id,
			Name:             prods.Name,
			Price:            prods.Price,
			Discount:         discount,
			Qty:              v.Quantity,
			TotalNormalPrice: prods.Price,
			TotalFinalPrice:  discCount,
		})

	}

	orderResp := models.Order{
		CashierID:      1,
		PaymentTypesId: body.PaymentTypeId,
		TotalPrice:     TotalInvoicePrice.ttprice,
		TotalPaid:      body.TotalPaid,
		TotalReturn:    body.TotalPaid - TotalInvoicePrice.ttprice,
		ReceiptId:      "R000" + strconv.Itoa(rand.Intn(1000)),
		ProductId:      productsIds,
		Quantities:     quantities,
		UpdatedAt:      time.Now().UTC(),
		CreatedAt:      time.Now().UTC(),
	}

	db.DB.Create(&orderResp)

	return c.Status(200).JSON(fiber.Map{

		"message": "success",
		"success": true,
		"data": map[string]interface{}{
			"order":    orderResp,
			"products": Prodresponse,
		},
	})

}
