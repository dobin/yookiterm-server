package main

import (
	"encoding/json"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

var HlSsoHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var token = r.FormValue("hlssotoken")
	isAuthenticated := false

	userId, err := getUsername(token)
	if err == nil {
		isAuthenticated = true
	} else {
		isAuthenticated = false
	}
	isAdmin := false

	logger.Infof("HL username: %s", userId)

	if isAuthenticated {
		token := userAuthToken(isAdmin, userId)
		t := "token=" + token + "; Path=/;"
		logger.Infof("Token: %s", t)
		w.Header().Set("Set-Cookie", t)

		http.Redirect(w, r, "//", http.StatusSeeOther)
		return
	} else {
		http.Redirect(w, r, "//", http.StatusSeeOther)
		return
	}
})

var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var userId string
	var password string

	type AuthReq struct {
		UserId   string
		Password string
	}

	decoder := json.NewDecoder(r.Body)
	var au AuthReq
	err := decoder.Decode(&au)
	if err != nil {
		//panic()
		fmt.Println("Error json")
	}

	userId = au.UserId
	password = au.Password
	isAdmin := false
	isAuthenticated := false

	// Authentication
	if userId == "admin" && password == config.AdminPassword {
		isAuthenticated = true
		isAdmin = true
		logger.Infof("Admin %s authenticated successfully", userId)
	} else if password == config.UserPassword {
		isAuthenticated = true
		isAdmin = false
		logger.Infof("User %s authenticated successfully", userId)
	}

	token := userAuthToken(isAdmin, userId)
	body := make(map[string]interface{})
	body["token"] = token

	if isAuthenticated {
		body["token"] = token
		body["authenticated"] = true
	} else {
		body["token"] = ""
		body["authenticated"] = false
	}

	err = json.NewEncoder(w).Encode(body)
	if err != nil {
		http.Error(w, "Internal server error", 500)
		return
	}
})

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

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return ([]byte(config.Jwtsecret)), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

var authTest = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Here we are converting the slice of products to json
	//payload, _ := json.Marshal("products")
	payload := "test"

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})
