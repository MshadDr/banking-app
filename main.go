package main

import (
	"banking_app/api"
	"banking_app/database"
	// "banking_app/migrations"
)

func main () {
	
	/* we did that once at the starting project
	 migrations.MigrateTransactions() */

	database.InitDatabase()
	api.StartApi()
}