package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/markbates/goth/gothic"
)

func userAuthToken(isAdmin bool, userId, userName string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin":    isAdmin,
		"userId":   userId,
		"userName": userName,
		"nbf":      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
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

		// The userId needs to be unique, so no collisions happening
		// e.g. for the created container.
		// It also should be of a constant length, so the hostname is
		// of a certain length, which makes the stack more predictable.
		// Part of hash of email address...
		// 4bit entropy * 8 chars = 32 bit entropy
		userHash := fmt.Sprintf("%x", sha256.Sum256([]byte(user.Email)))
		userId := string(userHash[:8])

		/*
			// Previously:
			// We would need some kind of "nickname"
			// But all we have is names and email
			// As it should be unique, just use all non-special letters from email
			reg, err := regexp.Compile("[^a-zA-Z0-9]+")
			if err != nil {
				log.Fatal(err)
			}
			userId := reg.ReplaceAllString(user.Email, "")
		*/

		// Auth success, set cookie
		token := userAuthToken(false, userId, user.Email)
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
