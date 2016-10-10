package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"math/rand"
	"time"

	"github.com/gorilla/mux"
	"github.com/lxc/lxd"
	"golang.org/x/exp/inotify"
	"gopkg.in/yaml.v2"
	"github.com/rs/cors"
	"github.com/howbazaar/loggo"
)


// Global variables
var lxdDaemon *lxd.Client
var config serverConfig

type serverConfig struct {
	ServerAddr          string   				`yaml:"server_addr"`
	ServerBannedIPs     []string 				`yaml:"server_banned_ips"`
	Jwtsecret						string   				`yaml:"jwtsecret"`
	ContainerHosts			[]ContainerHost	`yaml:"container_hosts"`
	BaseContainers			[]BaseContainer `yaml:"base_containers"`
}

type statusCode int

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


func parseConfig() error {
	data, err := ioutil.ReadFile("yookiterm-server.yml")
	if os.IsNotExist(err) {
		return fmt.Errorf("The configuration file (yookiterm-server.yml) doesn't exist.")
	} else if err != nil {
		return fmt.Errorf("Unable to read the configuration: %s", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("Unable to parse the configuration: %s", err)
	}

	if config.ServerAddr == "" {
		config.ServerAddr = ":8080"
	}

	return nil
}


func configWatcher() {
	// Watch for configuration changes
	watcher, err := inotify.NewWatcher()
	if err != nil {
		fmt.Errorf("Unable to setup inotify: %s", err)
	}

	err = watcher.Watch(".")
	if err != nil {
		fmt.Errorf("Unable to setup inotify watch: %s", err)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.Name != "./yookiterm-config.yml" {
					continue
				}

				if ev.Mask&inotify.IN_MODIFY != inotify.IN_MODIFY {
					continue
				}

				fmt.Printf("Reloading configuration\n")
				err := parseConfig()
				if err != nil {
					fmt.Printf("Failed to parse configuration: %s\n", err)
				}
			case err := <-watcher.Error:
				fmt.Printf("Inotify error: %s\n", err)
			}
		}
	}()
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
	r.Handle("/1.0/get-token", GetTokenHandler)
	r.Handle("/1.0/authTest", jwtMiddleware.Handler(authTest))

	r.Handle("/1.0/containerHosts", jwtMiddleware.Handler(restContainerHostListHandler))

	r.Handle("/1.0/baseContainers", jwtMiddleware.Handler(restBaseContainerListHandler))

	r.Handle("/1.0/challenges", jwtMiddleware.Handler(restChallengeListHandler))
	r.Handle("/1.0/challenge/{challengeId}", jwtMiddleware.Handler(restChallengeHandler))
	//r.HandleFunc("/1.0/challenge/<challenge>/file", restBaseContainerListHandler)

	c := cors.New(cors.Options{
	    AllowedOrigins: []string{"*"},
	    AllowCredentials: true,
			//Debug: true,
			AllowedHeaders: []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
	})
	handler := c.Handler(r)

	fmt.Println("Yookiterm server 0.2");
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
