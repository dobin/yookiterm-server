package main

import (
	"encoding/json"
//	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)


type BaseContainer struct {
	Name string
	Id string
}

type ContainerHost struct {
	HostnameAlias string
	Hostname string
}


var restBaseContainerListHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var baseContainerList [2]BaseContainer

	baseContainerList[0] = BaseContainer{"beta", "1"}
	baseContainerList[1] = BaseContainer{"hlUbuntu32aslr", "2"}

	err := json.NewEncoder(w).Encode(baseContainerList)
	if err != nil {
		http.Error(w, "Internal server error", 500)
		return
	}
})


var restContainerHostListHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(config.ContainerHosts)
	if err != nil {
		http.Error(w, "Internal server error", 500)
		return
	}
})


var restChallengeListHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var challengeList = getChallenges()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(challengeList)
	if err != nil {
		http.Error(w, "Internal server error", 500)
		return
	}
})


var restChallengeHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	challengeId := vars["challengeId"]

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	challenge := getChallenge(challengeId)

	err := json.NewEncoder(w).Encode(challenge)
	if err != nil {
		http.Error(w, "Internal server error", 500)
		return
	}
})
