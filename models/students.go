package models

import "gorm.io/gorm"

type Students struct {
	ID      uint   `gorm:"primary key;autoIncrement" json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Class   string `json:"class"`
	Subject string `json:"subject"`
}

func MigrateStudents(db *gorm.DB) error {
	err := db.AutoMigrate(&Students{})
	return err
}
