package models

import "time"

type Order struct {
	Id             int       `json:"id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	CashierID      int       `json:"cashier_id"`
	PaymentTypesId int       `json:"payment_types_id"`
	TotalPrice     int       `json:"total_price"`
	TotalPaid      int       `json:"total_paid"`
	TotalReturn    int       `json:"total_return"`
	ReceiptId      string    `json:"receipt_id"`
	IsDownload     int       `json:"is_download"`
	ProductId      string    `json:"product_id"`
	Quantities     string    `json:"quantities"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ProductResponseOrder struct {
	ProductId        int      `json:"productId"`
	Name             string   `json:"name"`
	Price            int      `json:"price"`
	Qty              int      `json:"qty"`
	Discount         Discount `json:"discount"`
	TotalNormalPrice int      `json:"totalNormalPrice"`
	TotalFinalPrice  int      `json:"totalFinalPrice"`
}

type ProductOrder struct {
	Id         int    `json:"id"`
	Sku        string `json:"sku"`
	Name       string `json:"name"`
	Stock      int    `json:"stock"`
	Price      int    `json:"price"`
	Image      string `json:"image"`
	CategoryId int    `json:"category_id"`
	DiscountId int    `json:"discount_id"`
}
