package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/howbazaar/loggo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/rs/cors"
)

var logger = loggo.GetLogger("project.main")

func main() {
	var err error

	rand.Seed(time.Now().UTC().UnixNano() + 0xcafebabe)

	err = run()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	var err error

	initLogger()
	fmt.Printf("Yookiterm server\n")

	// Setup configuration
	err = parseConfig()
	if err != nil {
		return err
	}

	// Load the challenges
	loadChallenges()

	// Setup the HTTP server
	r := mux.NewRouter()

	// Authentication Providers
	goth.UseProviders(
		google.New(config.GoogleId, config.GoogleSecret, config.ServerUrl+"/1.0/auth/google/callback", "email"),
	)

	// Authentication
	r.Handle("/1.0/auth/{provider}/callback", AuthProviderCallbackHandler)
	r.Handle("/1.0/auth/{provider}", AuthProviderHandler)
	r.Handle("/1.0/get-token", GetTokenHandler)

	// Authenticated
	r.Handle("/1.0/containerHosts", jwtMiddleware.Handler(restContainerHostListHandler))
	r.Handle("/1.0/baseContainers", jwtMiddleware.Handler(restBaseContainerListHandler))

	// Challenges are public
	r.Handle("/1.0/challenges", restChallengeListHandler)
	r.Handle("/1.0/challenge/{challengeId}", restChallengeHandler)
	//r.HandleFunc("/1.0/challenge/<challenge>/file", restBaseContainerListHandler)

	// Static Files: Slides
	r.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir(config.SlidesDir))))

	// Static HTML
	// Requires yookiterm project in parent directory
	// Should be at the end of the router
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(config.FrontendDir)))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		//Debug: true,
		AllowedHeaders: []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
	})
	handler := c.Handler(r)

	fmt.Printf("Listening on  : %s\n", config.ServerAddr)
	fmt.Printf("Serving domain: %s\n", config.ServerUrl)

	err = http.ListenAndServe(config.ServerAddr, handler)
	if err != nil {
		return err
	}

	return nil
}

func initLogger() {
	logger.SetLogLevel(loggo.DEBUG)
}
