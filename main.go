package main

import (
	"flag"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"./scenariolib"
	"github.com/k0kubun/pp"
)

func main() {
	// Init loggers

	tracePtr := flag.Bool("trace", false, "enable TRACE")
	seedPtr := flag.Int64("seed", -1, "set the Randomizer seed")

	flag.Parse()

	traceOut := ioutil.Discard
	if *tracePtr {
		pp.Println("TRACE enabled")
		traceOut = os.Stdout
	}

	scenariolib.InitLogger(traceOut, os.Stdout, os.Stdout, os.Stderr)

	seed := *seedPtr
	if seed == -1 {
		// Seed Random based on current time
		seed = int64(time.Now().Unix())
	}
	scenariolib.Trace.Printf("Ramdom seed: %d", seed)

	rand.Seed(seed)

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

	bot := scenariolib.NewUabot(local, scenarioURL, searchToken, analyticsToken)

	quit := make(chan bool)
	err := bot.Run(quit)
	if err != nil {
		scenariolib.Error.Println(err)
	}
	pp.Println("LOG >>> DONE")
}
