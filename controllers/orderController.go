package controllers

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	db "github.com/atifali-pm/go-pos-api/config"
	"github.com/atifali-pm/go-pos-api/middleware"
	"github.com/atifali-pm/go-pos-api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
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

	// Define a custom struct for Cashier without CreatedAt and UpdatedAt
	type CashierResponse struct {
		Id       uint   `json:"id"`
		Name     string `json:"name"`
		Passcode string `json:"passcode"`
	}

	type PaymentTypeResponse struct {
		Name string `json:"name"`
	}

	type CustomProductResponse struct {
		Sku      string          `json:"sku"`
		Name     string          `json:"name"`
		Stock    int             `json:"stock"`
		Price    int             `json:"price"`
		Image    string          `json:"image"`
		Category models.Category `json:"category"`
		Discount models.Discount `json:"discount"`
	}

	type OrderList struct {
		OrderId        int                     `json:"order_id"`
		CashierId      int                     `json:"cashierId"`
		PaymenttypesId int                     `json:"paymentTypesId"`
		TotalPrice     int                     `json:"totalPrice"`
		TotalPaid      int                     `json:"totalPaid"`
		TotalReturn    int                     `json:"totalReturn"`
		ReceiptId      string                  `json:"ReceiptId"`
		GeneratedAt    time.Time               `json:"generatedAt"`
		UpdatedAt      time.Time               `json:"updatedAt"`
		PaymentTypes   PaymentTypeResponse     `json:"payment_type"`
		Cashiers       CashierResponse         `json:"cashier"`
		Products       []CustomProductResponse `json:"products"`
	}

	OrderResponse := make([]*OrderList, 0)

	var categoryResult models.CategoryResult
	var discountResult models.DiscountResult

	for _, v := range order {
		cashier := models.Cashier{}
		db.DB.Where("id=?", v.CashierID).Find(&cashier)

		paymentType := models.PaymentType{}
		db.DB.Where("id=?", v.PaymentTypesId).Find(&paymentType)

		ProductIds := strings.Split(v.ProductId, ",")

		var products []models.Product

		// Define a slice to store product responses
		var productsResponse []CustomProductResponse

		for i := 1; i < len(ProductIds); i++ {
			product := models.Product{}
			db.DB.Where("id=?", ProductIds[i]).Find(&product)

			db.DB.Table("categories").Select("id, name").Where("id = ?", product.CategoryId).Find(&categoryResult)
			// Convert categoryResult to models.Category
			category := models.Category{
				Id:   categoryResult.Id,
				Name: categoryResult.Name,
			}

			db.DB.Table("discounts").Select("*").Where("id = ?", product.DiscountId).Limit(limit).Offset(skip).Find(&discountResult).Count(&count)
			// Convert categoryResult to models.Category
			discount := models.Discount{
				Id:   discountResult.Id,
				Type: discountResult.Type,
			}

			count = int64(len(products))

			products = append(products, product)

			// Populate the productsResponse slice
			productsResponse = append(productsResponse, CustomProductResponse{
				Sku:      product.Sku,
				Name:     product.Name,
				Stock:    product.Stock,
				Price:    product.Price,
				Image:    product.Image,
				Category: category,
				Discount: discount,
			})
		}

		OrderResponse = append(OrderResponse, &OrderList{
			OrderId:        v.Id,
			CashierId:      v.CashierID,
			PaymenttypesId: v.PaymentTypesId,
			TotalPrice:     v.TotalPrice,
			TotalPaid:      v.TotalPaid,
			TotalReturn:    v.TotalReturn,
			ReceiptId:      v.ReceiptId,
			GeneratedAt:    v.CreatedAt,
			PaymentTypes: PaymentTypeResponse{
				Name: paymentType.Name,
			},
			Cashiers: CashierResponse{
				Id:       cashier.Id,
				Name:     cashier.Name,
				Passcode: cashier.Passcode,
			},
			Products: productsResponse,
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

func DownloadOrder(c *fiber.Ctx) error {
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
			"Message": "Order Not Found",
			"error":   map[string]interface{}{},
		})
	}

	productIds := strings.Split(Order.ProductId, ",")

	TotalProducts := make([]*models.Product, 0)

	for i := 1; i < len(productIds); i++ {
		prods := models.Product{}
		db.DB.Where("id=?", productIds[i]).Find(&prods)
		TotalProducts = append(TotalProducts, &prods)
	}

	cashier := models.Cashier{}
	db.DB.Where("id = ?", Order.CashierID).Find(&cashier)

	paymentType := models.PaymentType{}
	db.DB.Where("id = ?", Order.PaymentTypesId).Find(&paymentType)

	orderTable := models.Order{}
	db.DB.Where("id = ?", Order.Id).Find(&orderTable)

	twoDarray := [][]string{{}}
	quantities := strings.Split(Order.Quantities, ",")
	quantities = quantities[1:]

	for i := 0; i < len(TotalProducts); i++ {
		s_array := []string{}
		s_array = append(s_array, TotalProducts[i].Sku)
		s_array = append(s_array, TotalProducts[i].Name)
		s_array = append(s_array, quantities[i])
		s_array = append(s_array, strconv.Itoa(TotalProducts[i].Price))
		twoDarray = append(twoDarray, s_array)
	}

	begin := time.Now()
	grayColor := getGrayColor()
	whiteColor := color.NewWhite()
	header := getHeader()
	contents := twoDarray

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)
	//m.SetBorder(true)

	//Top Heading
	m.SetBackgroundColor(grayColor)
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Order Invoice #"+strconv.Itoa(Order.Id), props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
	})
	m.SetBackgroundColor(whiteColor)

	//Table setting
	m.TableList(header, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{3, 4, 2, 3},
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{3, 4, 2, 3},
		},
		Align:                consts.Center,
		AlternatedBackground: &grayColor,
		HeaderContentSpace:   1,
		Line:                 false,
	})
	//Total price
	m.Row(20, func() {
		m.ColSpace(7)
		m.Col(2, func() {
			m.Text("Total:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text("RS. "+strconv.Itoa(Order.TotalPrice), props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Center,
			})
		})
	})
	m.Row(21, func() {
		m.ColSpace(7)
		m.Col(2, func() {
			m.Text("Total Paid:", props.Text{
				Top:   0.5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text("RS. "+strconv.Itoa(Order.TotalPaid), props.Text{
				Top:   0.5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Center,
			})
		})
	})

	m.Row(22, func() {
		m.ColSpace(7)
		m.Col(2, func() {
			m.Text("Total Return", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text("RS. "+strconv.Itoa(Order.TotalReturn), props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Center,
			})
		})
	})

	//Invoice creation
	currentTime := time.Now()
	pdfFileName := "invoice-" + currentTime.Format("2006-Jan-02")
	err := m.OutputFileAndClose(pdfFileName + ".pdf")
	if err != nil {
		fmt.Println("Could not save PDF:", err)
		os.Exit(1)
	}

	end := time.Now()
	fmt.Println(end.Sub(begin))

	//update recepit is downloaded to 1 means true
	db.DB.Table("orders").Where("id=?", Order.Id).Update("is_download", 1)
	return c.Status(200).JSON(fiber.Map{
		"Success": true,
		"Message": "Success",
	})

}

func getHeader() []string {
	return []string{"Product Sku", "Name", "Qty", "Price"}
}

func getGrayColor() color.Color {
	return color.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}
