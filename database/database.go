package database

import (
	"banking_app/helpers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	
	dburl := "host=postgres port=5432 user=postgres dbname=bankApp password=root sslmode=disable"
	// dburl := "postgres:://postgres:root@postgres:5432/bankApp"

	database, err := gorm.Open( postgres.Open( dburl ), &gorm.Config{} )
	helpers.HandleErr( err )

	DB = database

	// sqlDB, err := database.DB()
	// sqlDB.SetMaxIdleConns( 20 )
	// sqlDB.SetMaxOpenConns( 200 )
	// DB = sqlDB
}