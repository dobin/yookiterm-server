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

	// Setup configuration
	err = parseConfig()
	if err != nil {
		return err
	}

	// Start the configuration file watcher
	configWatcher()

	// Load the challenges
	loadChallenges()

	// Setup the HTTP server
	r := mux.NewRouter()

	// Authentication Providers
	clientId := ""
	clientSecret := ""
	goth.UseProviders(
		google.New(clientId, clientSecret, "http://exploit.courses/1.0/auth/google/callback", "email"),
	)

	// Authentication
	r.Handle("/1.0/auth/{provider}/callback", AuthProviderCallbackHandler)
	r.Handle("/1.0/auth/{provider}", AuthProviderHandler)

	// Authenticated
	r.Handle("/1.0/containerHosts", jwtMiddleware.Handler(restContainerHostListHandler))
	r.Handle("/1.0/baseContainers", jwtMiddleware.Handler(restBaseContainerListHandler))

	// Challenges are public
	r.Handle("/1.0/challenges", restChallengeListHandler)
	r.Handle("/1.0/challenge/{challengeId}", restChallengeHandler)
	//r.HandleFunc("/1.0/challenge/<challenge>/file", restBaseContainerListHandler)

	// Static HTML
	// Requires yookiterm project in parent directory
	// Should be at the end of the router
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../yookiterm/app/")))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		//Debug: true,
		AllowedHeaders: []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
	})
	handler := c.Handler(r)

	fmt.Println("Yookiterm server")
	fmt.Println("Listening on: ", config.ServerAddr)

	err = http.ListenAndServe(config.ServerAddr, handler)
	if err != nil {
		return err
	}

	return nil
}

func initLogger() {
	logger.SetLogLevel(loggo.DEBUG)
}
