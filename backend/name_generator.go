package main

import (
	"math/rand"
	"strings"
	"time"
)

const nameLengthMax = 12

func getRandomName() string {
	var randomName = ""
	for {
		randomName = strings.Title(randomAdjective()) + " " + randomAnimal()
		if len(randomName) < nameLengthMax {
			break
		}
	}
	return randomName
}

func randomAdjective() string {
	rand.Seed(time.Now().UTC().UnixNano())
	return adjectiveList[rand.Intn(cap(adjectiveList))]
}

func randomAnimal() string {
	rand.Seed(time.Now().UTC().UnixNano())
	return animalList[rand.Intn(cap(animalList))]
}
