package database

import (
	"firlydev/backend-test/features/applicant/data"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&data.Applicant{}); err != nil {
		return nil, err
	}
	return db, nil
}
