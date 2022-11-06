package api

import (
	"banking_app/helpers"
	"banking_app/interfaces"
	"banking_app/middlewares"
	"banking_app/controllers/transactions"
	"banking_app/controllers/useraccounts"
	"banking_app/controllers/users"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

// handle the login data from the router
type Login struct {
	Username string
	Password string 
}

// handle the register data from the router
type Register struct {
	Username string
	Password string
	Email string
}

// handle the transaction data from the body
type TransactionBody struct {
	UserId uint
	From uint
	To uint
	Amount int
}

// prepare the request body
func readBody( r *http.Request ) []byte {
	body, err := ioutil.ReadAll( r.Body )
	helpers.HandleErr( err )

	return body
}

func apiResponse( call map[string]interface{}, w http.ResponseWriter ) {
	// Prepare response...
	if call[ "message" ] == "all is fine..." {
		resp := call
		json.NewEncoder(w).Encode(resp)
	} else {
		// Handle error
		resp := call
		json.NewEncoder(w).Encode(resp)
	}
}

// LoginController
func login( w http.ResponseWriter, r *http.Request ) {
	//Ready body...
	body := readBody( r ) 

	//Handle Login...
	var formattedBody Login
	err := json.Unmarshal( body, &formattedBody )
	helpers.HandleErr( err )

	login := users.Login( formattedBody.Username, formattedBody.Password )

	apiResponse( login, w )
}

// RegisterController...
func register( w http.ResponseWriter, r *http.Request ) {
	//Ready body...
	body := readBody( r ) 

	//Handle register...
	var formattedBody Register
	err := json.Unmarshal( body, &formattedBody )
	helpers.HandleErr( err )

	register := users.Register( formattedBody.Username, formattedBody.Password, formattedBody.Email )

	apiResponse( register, w )

	// Prepare response...
	if register[ "message" ] == "all is fine..." {
		resp := register
		json.NewEncoder(w).Encode(resp)
	} else {
		// Handle error
		resp := interfaces.ErrResponse{ Message: "Invalid Inputs" }
		json.NewEncoder(w).Encode(resp)
	}
}

//
func GetUser( w http.ResponseWriter, r *http.Request ) {
	vars := mux.Vars( r )
	userId := vars[ "id" ]
	authToken := r.Header.Get( "Authorization" )

	user := users.GetUser( userId, authToken )
	apiResponse( user, w )
}

//
func GetMyTransactions( w http.ResponseWriter, r *http.Request ) {
	vars := mux.Vars( r )
	userId := vars[ "userId" ]
	authToken := r.Header.Get( "Authorization" )

	tranasactions := transactions.GetMyTransactions( userId, authToken )
	apiResponse( tranasactions, w )
}

//
func transaction( w http.ResponseWriter, r *http.Request ) {
	body := readBody( r )

	auth := r.Header.Get( "Authorization" )

	var formattedBody TransactionBody
	err := json.Unmarshal( body, &formattedBody )
	helpers.HandleErr( err )

	transaction := useraccounts.Transaction( formattedBody.UserId, formattedBody.From, formattedBody.To, formattedBody.Amount, auth )
	apiResponse( transaction, w )
}

// start router
func  StartApi() {
	router := mux.NewRouter()
	router.Use( helpers.PanicHandler )

	/* general routes */
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/register", register).Methods("POST")
	/* end of general routes */

	/* members routes */
	// authenticatedRouter := router.PathPrefix( "/member" ).Subrouter()
	// authenticatedRouter.Use( middlewares.CheckMember )
	// authenticatedRouter.HandleFunc("/transaction", transaction).Methods("POST")
	/* end of members routes */
	
	router.Handle( "/transaction", middlewares.CheckUserId( transaction ) ).Methods( "POST" )
	router.Handle("/user/{id}", middlewares.CheckMember( GetUser ) ).Methods("GET")
	router.Handle("/transaction/{userId}", middlewares.CheckMember( GetMyTransactions )).Methods("GET")

	fmt.Println("App is working on port 8080")
	const addr = "0.0.0.0:8088"
	server := http.Server{
		Handler: router,
		Addr: addr,
	}

	// Starting Server...
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal( "server failed" )
	}
}