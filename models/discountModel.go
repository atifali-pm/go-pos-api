package models

type Discount struct {
	Id int `json:"id" gorm:"type:INT(10) UNSINGED NOT NULL AUTO_INCREMENT;primaryKey"`
}
