package generation

import (
	"math/rand"
	"strings"
	"time"

	"github.com/castillobgr/sententia"
)

const nameLengthMax = 12

func GetRandomName() string {
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

func GetRandomMessage() string {
	rand.Seed(time.Now().UTC().UnixNano())
	template := templateList[rand.Intn(cap(templateList))]

	sentence, err := sententia.Make(template)
	if err != nil {
		panic(err)
	}
	return sentence
}
