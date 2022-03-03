package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

var challengeList []Challenge

type Challenge struct {
	Id                 string
	ContainerBaseName  string
	ContainerHostAlias string

	Title       string
	Description string

	TextFilename string
	Text         string
	Active       bool

	BaseContainer sBaseContainer
	ContainerHost sContainerHost
}

// Load all challenges from the challenge directory
func loadChallenges() {
	challengesDir := config.ChallengesDir

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

			challenge.BaseContainer = getBaseContainerByName(challenge.ContainerBaseName)
			challenge.ContainerHost = getContainerHostByAlias(challenge.ContainerHostAlias)

			if challenge.Active {
				challengeList = append(challengeList, challenge)
			}
		}
	}
	fmt.Printf("Loaded %d challenges\n", len(challengeList))
}

func getChallenges() []Challenge {
	var strippedChallenges []Challenge

	// Remove challenge texts
	for _, element := range challengeList {
		challenge := element // copy challenge
		challenge.Text = ""  // remove text of challenge
		strippedChallenges = append(strippedChallenges, challenge)
	}

	return strippedChallenges
}

func getChallenge(challengeId string) (error, Challenge) {
	for index, element := range challengeList {
		if challengeList[index].Id == challengeId {
			return nil, element
		}
	}

	return errors.New("Challenge does not exist"), Challenge{}
}
