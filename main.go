package main

import (
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/coveo/uabot/scenariolib"
	"github.com/k0kubun/pp"
)

func main() {
	// Init loggers
	scenariolib.InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	// Seed Random based on current time
	source := rand.NewSource(int64(time.Now().Unix()))
	random := rand.New(source)

	searchToken := os.Getenv("SEARCHTOKEN")
	analyticsToken := os.Getenv("UATOKEN")
	if searchToken == "" || analyticsToken == "" {
		scenariolib.Error.Println("SEARCHTOKEN, UATOKEN need to be defined as env variables")
	}

	scenarioURL := os.Getenv("SCENARIOSURL")
	if scenarioURL == "" {
		scenariolib.Error.Println("SCENARIOSURL env variable needs to define a file path")
	}

	local := os.Getenv("LOCAL") == "true"
	if local {
		scenariolib.Info.Println("STARTING IN LOCAL MODE, MAKE SURE THE SCENARIOSURL IS A LOCAL PATH")
	}

	bot := scenariolib.NewUabot(local, scenarioURL, searchToken, analyticsToken, random)

    quit := make(chan bool)
	err := bot.Run(quit)
	if err != nil {
		scenariolib.Error.Println(err)
	}
	pp.Println("LOG >>> DONE")
}
