package controllers

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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
		CashierId     int        `json:"cashier_id"`
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
		CashierID:      body.CashierId,
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

func SubTotalOrder(c *fiber.Ctx) error {
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
		Products []products `json:"products"`
	}{}

	if err := c.BodyParser(&body.Products); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Empty Body",
		})
	}

	Prodresponse := make([]*models.ProductResponseOrder, 0)

	var TotalInvoicePrice = struct {
		ttprice int
	}{}

	for _, v := range body.Products {
		totalPrice := 0

		prods := models.ProductOrder{}

		var discount models.Discount
		db.DB.Table("products").Where("id=?", v.ProductId).First(&prods)
		db.DB.Where("id=?", prods.DiscountId).Find(&discount)

		disc := 0
		if discount.Type == "PERCENT" {
			totalPrice = prods.Price * v.Quantity
			percentage := totalPrice * discount.Result / 100
			disc = totalPrice - percentage
		}

		if discount.Type == "BUY_IN" {
			totalPrice = prods.Price * v.Quantity
			disc = totalPrice - discount.Result
		}

		TotalInvoicePrice.ttprice = TotalInvoicePrice.ttprice + disc

		Prodresponse = append(Prodresponse, &models.ProductResponseOrder{
			ProductId:        prods.Id,
			Name:             prods.Name,
			Price:            prods.Price,
			Discount:         discount,
			Qty:              v.Quantity,
			TotalNormalPrice: prods.Price,
			TotalFinalPrice:  disc,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
		"success": true,
		"data": map[string]interface{}{
			"subTotal": TotalInvoicePrice.ttprice,
			"products": Prodresponse,
		},
	})

}

func CheckOrder(c *fiber.Ctx) error {
	// Token authenticate
	if err := middleware.AuthorizeToken(c); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	// Token authenticate

	OrderId := c.Params("order_id")

	var order models.Order
	db.DB.Where("id=?", OrderId).First(&order)

	if order.Id <= 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": true,
			"message": "Order id not found!",
		})
	}

	if order.IsDownload <= 0 {
		return c.Status(200).JSON(fiber.Map{
			"status":  true,
			"message": "success",
			"data": map[string]interface{}{
				"isDownloaded": false,
			},
		})
	}

	if order.IsDownload > 0 {
		return c.Status(200).JSON(fiber.Map{
			"status":  true,
			"message": "success",
			"data": map[string]interface{}{
				"isDownloaded": true,
			},
		})
	}

	return nil
}

func OrderDetail(c *fiber.Ctx) error {
	// Token authenticate
	if err := middleware.AuthorizeToken(c); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	// Token authenticate

	OrderId := c.Params("order_id")
	var Order models.Order

	db.DB.Where("id=?", OrderId).First(&Order)

	if Order.Id <= 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Order not found!",
			"error":   map[string]interface{}{},
		})
	}

	ProductIds := strings.Split(Order.ProductId, ",")
	TotalProducts := make([]*models.Product, 0)

	for i := 1; i < len(ProductIds); i++ {
		prods := models.Product{}
		db.DB.Where("id=?", ProductIds[i]).Find(&prods)
		TotalProducts = append(TotalProducts, &prods)

	}

	cashier := models.Cashier{}
	db.DB.Where("id=?", Order.CashierID).Find(&cashier)
	fmt.Println(cashier)

	paymentType := models.PaymentType{}
	db.DB.Where("id=?", Order.PaymentTypesId).Find(&paymentType)

	orderTable := models.Order{}
	db.DB.Where("id=?", Order.Id).Find(&orderTable)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"data": map[string]interface{}{
			"order": map[string]interface{}{
				"OrderId":        Order.Id,
				"CashierId":      Order.CashierID,
				"PaymentTypesID": Order.PaymentTypesId,
				"TotalPrice":     Order.TotalPrice,
				"TotalPaid":      Order.TotalPaid,
				"ReceiptId":      Order.ReceiptId,
				"CreatedAt":      Order.CreatedAt,
				"Cashier":        cashier,
				"PaymentType":    paymentType,
			},
			"products": TotalProducts,
		},
		"message": "success",
	})
}

func OrderList(c *fiber.Ctx) error {
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
	var order []models.Order
	db.DB.Select("*").Limit(limit).Offset(skip).Find(&order).Count(&count)

	type OrderList struct {
		OrderId        int                `json:"order_id"`
		CashierId      int                `json:"cashierId"`
		PaymenttypesId int                `json:"paymentTypesId"`
		TotalPrice     int                `json:"totalPrice"`
		TotalPaid      int                `json:"totalPaid"`
		TotalReturn    int                `json:"totalReturn"`
		ReceiptId      string             `json:"ReceiptId"`
		GeneratedAt    time.Time          `json:"generatedAt"`
		UpdatedAt      time.Time          `json:"updatedAt"`
		PaymentTypes   models.PaymentType `json:"paymentType"`
		Cashiers       models.Cashier     `json:"cashier"`
	}

	OrderResponse := make([]*OrderList, 0)

	for _, v := range order {

		cashier := models.Cashier{}
		db.DB.Where("id=?", v.CashierID).Find(&cashier)

		paymentType := models.PaymentType{}
		db.DB.Where("id=?", v.PaymentTypesId).Find(&paymentType)

		OrderResponse = append(OrderResponse, &OrderList{
			OrderId:        v.Id,
			CashierId:      v.CashierID,
			PaymenttypesId: v.PaymentTypesId,
			TotalPrice:     v.TotalPrice,
			TotalPaid:      v.TotalPaid,
			TotalReturn:    v.TotalReturn,
			ReceiptId:      v.ReceiptId,
			GeneratedAt:    v.CreatedAt,
			PaymentTypes:   paymentType,
			Cashiers:       cashier,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data": map[string]interface{}{
			"orders": OrderResponse,
		},
		"meta": map[string]interface{}{
			"total": count,
			"limit": limit,
			"skip":  skip,
		},
	})

}
