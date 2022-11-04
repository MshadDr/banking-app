package migrations

import (
	"banking_app/database"
	"banking_app/helpers"
	"banking_app/interfaces"
)

func createAccounts() {
	users := &[2]interfaces.User{
		{Username: "Martin", Email: "martin@gmail.com"},
		{Username: "Sam", Email: "sam@gmail.com"},
	}

	for i := 0; i < len(users); i++ {
		
		generatedPassword := helpers.HashAndSalt( []byte( users[ i ].Username ) )
		user := &interfaces.User{ Username: users[i].Username, Email: users[i].Email, Password: generatedPassword }

		database.DB.Create( &user )

		account := &interfaces.Account{Type: "Daily Account", Name: string( users[i].Username + "'s" + " account" ), Balance: uint( 10000 * int( i + 1 ) ), UserID: user.ID }

		database.DB.Create( &account )
	}
}

// migrate and create tables
func Migrate () {

	User := &interfaces.User{}
	Account := &interfaces.Account{}
	Transaction := &interfaces.Transaction{}

	database.DB.AutoMigrate( User, Account, Transaction )

	createAccounts()
}