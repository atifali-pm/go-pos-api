package models

import "time"

type Order struct {
	Id             int       `json:"id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	CashierID      int       `json:"cashier_id"`
	PaymentTypesID int       `json:"payment_types_id"`
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
