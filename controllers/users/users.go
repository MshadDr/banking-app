package users

import (
	"banking_app/database"
	"banking_app/helpers"
	"banking_app/interfaces"
	"time"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// function for generate a jwt token for @user
func prepareToken( user *interfaces.User ) string {
	
	tokenContent := jwt.MapClaims{
		"user_id": user.ID,
		"expiry": time.Now().Add(time.Minute ^ 60).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleErr(err)

	return token
}

// function for Setup @response
func prepareResponse( user *interfaces.User, accounts []interfaces.ResponseAccount, withToken bool ) map[string]interface{} {

	responseUser := &interfaces.ResponseUser{
		ID: user.ID,
		Username: user.Username,
		Email: user.Email,
		Account: accounts,
	}

	var response = map[string]interface{}{ "message": "all is fine..." }
	if withToken {
		var token = prepareToken( user )
		response[ "jwt" ] = token
	}
	
	response["data"] = responseUser

	return response

}

// func for login By @username && @password
func Login( username string, password string ) map[string]interface{} {

	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: password, Valid: "password"},
		})
		if valid {
			return successValidation( username, password )

		} else {
			return map[ string ]interface{}{ "message": "not valid values" }
		}
}

// func for Register new user by @username && @email && @password
func Register( username string, password string, email string ) map[string]interface{} {

		valid := helpers.Validation(
			[]interfaces.Validation{
				{Value: username, Valid: "username" },
				{Value: password, Valid: "password" },
				{Value: email, Valid: "email" },
			})
		if valid {
			return createUser( username, password, email )
		} else {
			 return map[ string ]interface{}{ "message": "not valid value!!!" }
		}
}

// after validation go to login...
func successValidation( username string, password string ) map[string]interface{} {

	user := &interfaces.User{}
	if err := database.DB.Where( "username = ?", username).First(&user).Error; err != nil {
		return map[string]interface{}{"message": "User not found..."}
	}

	// Verify Password
	passErr := bcrypt.CompareHashAndPassword( []byte(user.Password), []byte( password ) )

	if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
		return map[string]interface{}{"message": "Wrong password..."}
	}

	// Find account for the user...
	accounts := []interfaces.ResponseAccount{}
	database.DB.Table("accounts").Select("id, name, balance").Where("user_id = ?", user.ID).Scan(&accounts)

	// Prepare response...
	return prepareResponse( user, accounts, true )

}

// create a new user process...
func createUser( username string, password string, email string ) map[string]interface{} {

	generatedPassword := helpers.HashAndSalt( []byte(password) )

	user := &interfaces.User{ Username: username, Email: email, Password: generatedPassword }
	database.DB.Create( &user )

	account := &interfaces.Account{Type: "Daily Account", Name: string( username + "'s" + " account" ), Balance: 0, UserID: user.ID }
	database.DB.Create( &account )

	accounts := []interfaces.ResponseAccount{}
	respAccount := interfaces.ResponseAccount{ ID: account.ID, Name: account.Name, Balance: int( account.Balance )}
	accounts = append( accounts, respAccount )
	var response = prepareResponse( user, accounts, true )
	return response;
}

func GetUser( id string, jwt string) map[string]interface{} {
	isValid := helpers.ValidateToken( id, jwt )

	if isValid {

		user := &interfaces.User{}
		if err := database.DB.Where( "id = ?", id).First(&user).Error; err != nil {
			return map[string]interface{}{"message": "User not found..."}
		}

		// Find account for the user...
		accounts := []interfaces.ResponseAccount{}
		database.DB.Table("accounts").Select( "id, name, balance" ).Where( "user_id = ?", id ).Scan(&accounts)

		var response = prepareResponse( user, accounts, false )
		return response
	} else {
		return map[string]interface{}{ "message": "Not valid token" }
	}
}