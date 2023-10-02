package models

import (
	"encoding/json"
	"time"
)

type Cashier struct {
	Id        uint      `json:"id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Name      string    `json:"name"`
	Passcode  string    `json:"passcode"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CashierResponse struct {
	Cashier
}

func (c CashierResponse) MarshalJSON() ([]byte, error) {
	type Alias CashierResponse
	return json.Marshal(&struct {
		CreatedAt string `json:"created_at,omitempty"`
		UpdatedAt string `json:"updated_at,omitempty"`
		*Alias
	}{
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
		Alias:     (*Alias)(&c),
	})
}
