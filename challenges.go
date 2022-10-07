package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

var challengeList []Challenge

type Challenge struct {
	Id                 string
	ContainerBaseName  string
	ContainerHostAlias string

	Title       string
	Description string
	Chapter		string

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
		// read json file in challenges/challengeX/
		filename := challengesDir + "/" + d.Name() + "/" + d.Name() + ".json"

		fileContent, e := ioutil.ReadFile(filename)
		if e != nil {
			fmt.Println("Error reading file: ", filename)
		}

		// create challenge based on json file
		var challenge Challenge
		json.Unmarshal(fileContent, &challenge)

		// also read markup file, as referenced in the json file
		var markupFilename = challengesDir + "/" + d.Name() + "/" + d.Name() + ".md"
		markupFileContent, e := ioutil.ReadFile(markupFilename)
		if e != nil {
			fmt.Println("Error reading file: ", filename)
		}
		challenge.Text = string(markupFileContent)

		challenge.Id = d.Name()[len(d.Name())-2:]  // last two characters are the id

		challenge.BaseContainer = getBaseContainerByName(challenge.ContainerBaseName)
		challenge.ContainerHost = getContainerHostByAlias(challenge.ContainerHostAlias)

		if challenge.Active {
			challengeList = append(challengeList, challenge)
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
