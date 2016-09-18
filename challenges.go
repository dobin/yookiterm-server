package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
)

var challengeList []Challenge

var challengesDir = "./challenges"


type Challenge struct {
	Id string
	ContainerBaseName string
	ContainerArch int
	ContainerAslr string
	ContainerUser string
	Title string
	Text string
	TextFilename string
	ContainerHostAlias string
}


// Load all challenges from the challenge directory
func loadChallenges() {
	fmt.Println("Loading challenges")

	// read all directories in challenges directory
	// ./challenges/*
	dirs, _ := ioutil.ReadDir(challengesDir)
	for _, d := range dirs {

		// read json files in challenges/challengX/
		files, _ := filepath.Glob(challengesDir + "/" + d.Name() + "/" + "*.json")
		for _, f := range files {
			// read json file
			fileContent, e := ioutil.ReadFile(f)
			if e != nil {
				fmt.Println("Error reading file: ", f)
			}

			// create challenge based on json file
			var challenge Challenge
			json.Unmarshal(fileContent, &challenge)

			// also read markup file, as referenced in the json file
			var markupFilename = challenge.TextFilename
			markupFileContent, ee := ioutil.ReadFile(challengesDir + "/" + d.Name() + "/" + markupFilename)
			if ee != nil {
				fmt.Println("Error reading file: ", f)
			}
			challenge.Text = string(markupFileContent)

			//if config.ContainerDomain != "" {
			//	challenge.ContainerHostFQDN = challenge.ContainerHost + "." + config.ContainerDomain
			//}

			challengeList = append(challengeList, challenge)
		}
	}
}


func getChallenges() []Challenge {
	return challengeList
}


func getChallenge(challengeId string) Challenge {

	for index, element := range challengeList {
		if challengeList[index].Id == challengeId {
			return element
		}
	}

	// TODO FIXME ret error
	return challengeList[0]
}
