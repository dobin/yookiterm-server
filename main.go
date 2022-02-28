package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/howbazaar/loggo"
	"github.com/rs/cors"
)

// Global variables
var logger = loggo.GetLogger("project.main")

func main() {
	rand.Seed(time.Now().UTC().UnixNano() + 0xcafebabe)

	var err error

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

	// Authentication
	r.Handle("/1.0/hl-sso-auth", HlSsoHandler)
	r.Handle("/1.0/get-token", GetTokenHandler)
	r.Handle("/1.0/authTest", jwtMiddleware.Handler(authTest))

	r.Handle("/1.0/containerHosts", jwtMiddleware.Handler(restContainerHostListHandler))
	r.Handle("/1.0/baseContainers", jwtMiddleware.Handler(restBaseContainerListHandler))

	// Challenges are public
	r.Handle("/1.0/challenges", restChallengeListHandler)
	r.Handle("/1.0/challenge/{challengeId}", restChallengeHandler)
	//r.HandleFunc("/1.0/challenge/<challenge>/file", restBaseContainerListHandler)

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
