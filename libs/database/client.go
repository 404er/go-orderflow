package database

import (
	"log"
	orderflow "orderFlow/libs/orderflow"
	"orderFlow/libs/shared"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabaseClient() {
	db, err := gorm.Open(postgres.Open(shared.DB_URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&orderflow.FootprintCandle{})
	if err != nil {
		panic(err)
	}
	log.Println("Client database success")
	DB = db
}
