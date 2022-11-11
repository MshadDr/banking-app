package database

import (
	"banking_app/helpers"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {

	// database, err := gorm.Open( "postgres", "host=localhost port=5432 user=postgres dbname=bankApp password=root sslmode=disable" )
	fmt.Println(1)
	dburl := "postgres:://postgres:root@postgres:5432/bankApp"

	database, err := gorm.Open( postgres.Open( dburl ), &gorm.Config{} )
	fmt.Println(2)
	helpers.HandleErr( err )

	DB = database

	// sqlDB, err := database.DB()

	// sqlDB.SetMaxIdleConns( 20 )
	// sqlDB.SetMaxOpenConns( 200 )
	// DB = sqlDB
}