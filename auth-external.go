package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
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

		// We would need some kind of "nickname"
		// But all we have is names and email
		// As it should be unique, just use all non-special letters from email
		reg, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			log.Fatal(err)
		}
		userId := reg.ReplaceAllString(user.Email, "")

		// Auth success, set cookie
		token := userAuthToken(false, userId)
		t := "token=" + token + "; Path=/;"
		logger.Infof("Token: %s", t)
		res.Header().Set("Set-Cookie", t)

		// Redirect
		http.Redirect(res, req, "/", http.StatusPermanentRedirect)
	},
)

var AuthProviderHandler = http.HandlerFunc(
	func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	},
)
