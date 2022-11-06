package middlewares

import (
	"net/http" 
	"strings"
	"banking_app/helpers"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"strconv"
	"errors"
)

func CheckMember( next http.HandlerFunc ) http.HandlerFunc {
	return http.HandlerFunc( func ( w http.ResponseWriter, r *http.Request ) {

		jwtToken := r.Header.Get( "Authorization" )
		cleanJWT := strings.Replace( jwtToken, "Bearer ", "", -1 )
		tokenData := jwt.MapClaims{}

		_ ,err := jwt.ParseWithClaims( cleanJWT, tokenData, func ( token *jwt.Token ) ( interface{}, error ) {
			return []byte( "TokenPassword" ), nil
		} ) 
		
		helpers.HandleErr( err )
		next.ServeHTTP( w, r )
	} )
}

// check the user id with JWT token
func CheckUserId( next http.HandlerFunc ) http.HandlerFunc {
	return http.HandlerFunc( func ( w http.ResponseWriter, r *http.Request ) {
		jwtToken := r.Header.Get( "Authorization" )
		cleanJWT := strings.Replace( jwtToken, "Bearer ", "", -1 )
		tokenData := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims( cleanJWT, tokenData, func ( token *jwt.Token ) ( interface{}, error ) {
			return []byte( "TokenPassword" ), nil
		} ) 
		helpers.HandleErr( err )

		vars := mux.Vars( r )
		id := vars[ "userId" ]

		var userId, _ = strconv.ParseFloat( id, 8 )

		if  token.Valid && tokenData[ "user_id" ] == userId {
			next.ServeHTTP( w, r )
		 }

		helpers.HandleErr( errors.New(" token is not valid ") )
	})
}