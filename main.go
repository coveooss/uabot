package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/erocheleau/uabot/scenariolib"
	"github.com/k0kubun/pp"
)

const (
	// USERAGENT This is the user agent the bot appears to be using.
	USERAGENT string = "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.80 Safari/537.36"
	// TIMEBETWEENVISITS The time for the bot to wait between visits, between 0 and X Seconds
	TIMEBETWEENVISITS int = 120
)

func main() {
	rand.Seed(int64(time.Now().Unix()))

	searchToken := os.Getenv("SEARCHTOKEN")
	analyticsToken := os.Getenv("UATOKEN")
	if searchToken == "" || analyticsToken == "" {
		pp.Fatal("FATAL >>> SEARCHTOKEN, UATOKEN need to be defined as env variables")
	}

	scenarioURL := os.Getenv("SCENARIOSURL")
	if scenarioURL == "" {
		pp.Fatal("FATAL >>> SCENARIOSURL env variable needs to define a file path")
	}

	timeNow := time.Now()

	// Init from path instead of URL, for testing purposes
	//conf, err := config.NewConfigFromPath(scenarioURL)
	conf, err := scenariolib.NewConfigFromURL(scenarioURL)
	if err != nil {
		pp.Fatal(err)
		return
	}

	count := 0
	for { // Run forever

		// Refresh the scenario files every 5 hours automatically.
		// This way, no need to stop the bot to update the possible scenarios.
		if time.Since(timeNow).Hours() > 5 {
			pp.Println("LOG >>> Updating Scenario file")
			// Init from path instead of URL, for testing purposes
			//conf, err := config.NewConfigFromPath(scenarioURL)
			conf2, err := scenariolib.NewConfigFromURL(scenarioURL)
			if err != nil {
				pp.Println("WARN >>> Cannot update scenario file, keeping the old one")
			} else {
				conf = conf2
			}
			timeNow = time.Now()
		}

		// New visit
		visit, err := scenariolib.NewVisit(searchToken, analyticsToken, USERAGENT, conf)
		if err != nil {
			pp.Fatal(err)
			return
		}

		// Setup specific stuff for NTO
		visit.SetupNTO()

		err = visit.ExecuteRandomScenario(conf)
		if err != nil {
			pp.Fatal(err)
			return
		}

		visit.UAClient.DeleteVisit()
		time.Sleep(time.Duration(rand.Intn(TIMEBETWEENVISITS)) * time.Second)

		count++
		pp.Printf("\n%v scenarios executed\n", count)
	}
	pp.Println("LOG >>> DONE")
}
