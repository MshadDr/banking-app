package useraccounts

import (
	"banking_app/helpers"
	"banking_app/database"
	"banking_app/interfaces"
	"banking_app/controllers/transactions"
	"strconv"
)

// update an account
func updateAccount( id uint, amount int ) interfaces.ResponseAccount {

	account := interfaces.Account{}
	responseAcc := interfaces.ResponseAccount{}

	database.DB.Where( "id = ? ", id ).First( &account )
	account.Balance = uint( amount )
	database.DB.Save( &account )

	responseAcc.ID = account.ID
	responseAcc.Name = account.Name
	responseAcc.Balance = int(account.Balance)

	return responseAcc
}

// get an account
func getAccount ( id uint ) *interfaces.Account {
	account := &interfaces.Account{}
	if database.DB.Where( "id = ? ", id ).First( &account ).RecordNotFound() {
		return nil
	}

	return account
}

// submit a transaction
func Transaction ( userId uint, from uint, to uint, amount int, jwt string ) map[string]interface{} {
	userIdString := strconv.FormatUint(uint64(userId), 10)
	isValid := helpers.ValidateToken( userIdString, jwt )

	if isValid {
		fromAccount := getAccount( from )
		toAccount := getAccount( to )

		if fromAccount == nil || toAccount == nil {
			return map[string]interface{}{ "message": "Account not found" }
		} else if fromAccount.UserID != userId {
			return map[string]interface{}{ "message": "Account not found" }
		} else if int( fromAccount.Balance ) < amount {
			return map[string]interface{}{ "message": "Account balance is too small" }
		}

		updatedAccount := updateAccount( from, int( fromAccount.Balance ) - amount )
		updateAccount( to, int( toAccount.Balance ) + amount )

		transactions.CreateTransaction( from, to, amount )

		var response = map[string]interface{}{ "message": "all is fine" }
		response[ "data" ] = updatedAccount

		return response
	} else {
		return map[string]interface{}{ "message": "Not valid token" }
	}
}