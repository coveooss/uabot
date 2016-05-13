package main

import (
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/erocheleau/uabot/scenariolib"
	"github.com/k0kubun/pp"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

const (
	// USERAGENT This is the user agent the bot appears to be using.
	//USERAGENT string = "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.80 Safari/537.36"

	// TIMEBETWEENVISITS The time for the bot to wait between visits, between 0 and X Seconds
	TIMEBETWEENVISITS int = 1
)

func Init(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {

	Trace = log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime)

	Warning = log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	// Init loggers
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	// Seed Random based on current time
	rand.Seed(int64(time.Now().Unix()))

	searchToken := os.Getenv("SEARCHTOKEN")
	analyticsToken := os.Getenv("UATOKEN")
	if searchToken == "" || analyticsToken == "" {
		Error.Println("SEARCHTOKEN, UATOKEN need to be defined as env variables")
	}

	scenarioURL := os.Getenv("SCENARIOSURL")
	if scenarioURL == "" {
		Error.Println("SCENARIOSURL env variable needs to define a file path")
	}

	local := os.Getenv("LOCAL")
	if local == "true" {
		Info.Println("STARTING IN LOCAL MODE, MAKE SURE THE SCENARIOSURL IS A LOCAL PATH")
	}

	timeNow := time.Now()

	var conf *scenariolib.Config
	var err error
	// Init from path instead of URL, for testing purposes
	if local == "true" {
		conf, err = scenariolib.NewConfigFromPath(scenarioURL)
	} else {
		conf, err = scenariolib.NewConfigFromURL(scenarioURL)
	}
	//
	if err != nil {
		Error.Println(err)
		return
	}

	count := 0
	for { // Run forever

		// Refresh the scenario files every 5 hours automatically.
		// This way, no need to stop the bot to update the possible scenarios.
		if time.Since(timeNow).Hours() > 5 {
			// false for local filepath, true for web hosted file
			var isURL = true
			if local == "true" {
				isURL = false
			}
			conf2 := refreshScenarios(scenarioURL, isURL)
			if conf2 != nil {
				conf = conf2
			}
			timeNow = time.Now()
		}

		scenario, err := conf.RandomScenario()
		if err != nil {
			Error.Println(err)
		}

		if scenario.UserAgent == "" {
			scenario.UserAgent, err = conf.RandomUserAgent(false)
			if err != nil {
				Error.Println(err)
			}
		}

		// New visit
		visit, err := scenariolib.NewVisit(searchToken, analyticsToken, scenario.UserAgent, conf)
		if err != nil {
			Error.Println(err)
			return
		}

		// Setup specific stuff for NTO
		//visit.SetupNTO()
		// Use this line instead outside of NTO
		visit.SetupGeneral()
		visit.LastQuery.CQ = conf.GlobalFilter

		err = visit.ExecuteScenario(*scenario, conf)
		if err != nil {
			Error.Println(err)
		}

		visit.UAClient.DeleteVisit()
		time.Sleep(time.Duration(rand.Intn(TIMEBETWEENVISITS)) * time.Second)

		count++
		Info.Printf("Scenarios executed : %d \n =============================\n\n", count)

	}
	pp.Println("LOG >>> DONE")
}

func refreshScenarios(url string, isUrl bool) *scenariolib.Config {
	Info.Println("Updating Scenario file")

	var err error
	var conf *scenariolib.Config

	if isUrl {
		conf, err = scenariolib.NewConfigFromURL(url)
	} else {
		conf, err = scenariolib.NewConfigFromPath(url)
	}

	if err != nil {
		Warning.Println("Cannot update scenario file, keeping the old one")
		return nil
	} else {
		return conf
	}
}
