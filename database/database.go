package database

import (
	"banking_app/helpers"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDatabase() {

	database, err := gorm.Open( "postgres", "host=localhost port=5432 user=postgres dbname=bankApp password=root sslmode=disable" )
	helpers.HandleErr( err )
	
	database.DB().SetMaxIdleConns( 20 )
	database.DB().SetMaxOpenConns( 200 )
	DB = database
}