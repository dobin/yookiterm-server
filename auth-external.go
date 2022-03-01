package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/markbates/goth/gothic"
)

func userAuthToken(isAdmin bool, userId string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin":  isAdmin,
		"userId": userId,
		"nbf":    time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(config.Jwtsecret))

	return tokenString
}

var AuthProviderCallbackHandler = http.HandlerFunc(
	func(res http.ResponseWriter, req *http.Request) {
		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}

		// Auth success, set cookie
		token := userAuthToken(false, user.Email)
		t := "token=" + token + "; Path=/;"
		logger.Infof("Token: %s", t)
		res.Header().Set("Set-Cookie", t)
	},
)

var AuthProviderHandler = http.HandlerFunc(
	func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	},
)
