package services

import (
	"fmt"
	"net/http"

	_ "github.com/dgrijalva/jwt-go"
)

var i int = 0

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i++

		fmt.Printf("This is a middleware (%v). Request #%v\n", r.RequestURI, i)

		next.ServeHTTP(w, r)
	})
}
