package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	Id        int       `json:"id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Category) BeforeSave(tx *gorm.DB) (err error) {
	// This hook is called before the model is saved
	// You can perform any migration-related tasks here

	// Check if the table exists
	tableExists := tx.Migrator().HasTable(&Category{})

	if !tableExists {
		// If the table doesn't exist, create the table
		if err := tx.AutoMigrate(&Category{}); err != nil {
			return err
		}
	}
	return
}
