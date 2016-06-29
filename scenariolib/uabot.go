package scenariolib

import (
	"time"
	"math/rand"
)

// DEFAULTTIMEBETWEENVISITS The time for the bot to wait between visits, between 0 and X Seconds
const DEFAULTTIMEBETWEENVISITS int = 120

type Uabot interface {
	Run() error
}

type uabot struct {
	local          bool
	conf           *Config
	scenarioURL    string
	searchToken    string
	analyticsToken string
	random         *rand.Rand
}

func NewUabot(local bool, conf *Config, scenarioUrl string, searchToken string, analyticsToken string, random *rand.Rand) *uabot {
	return &uabot{
		local,
		conf,
		scenarioUrl,
		searchToken,
		analyticsToken,
		random,
	}
}

func (bot *uabot) Run() error {
	// false for local filepath, true for web hosted file
	var isURL = true
	if bot.local {
		isURL = false
	}

	count := 0
	timeNow := time.Now()
	for { 	// Run forever
		// Refresh the scenario files every 5 hours automatically.
		// This way, no need to stop the bot to update the possible scenarios.
		if time.Since(timeNow).Hours() > 5 {
			conf2 := refreshScenarios(bot.scenarioURL, isURL)
			if conf2 != nil {
				bot.conf = conf2
			}
			timeNow = time.Now()
		}
		var timeVisits int
		if bot.conf.TimeBetweenVisits > 0 {
			timeVisits = bot.conf.TimeBetweenVisits
		} else {
			timeVisits = DEFAULTTIMEBETWEENVISITS
		}

		scenario, err := bot.conf.RandomScenario()
		if err != nil {
			return err
		}

		if scenario.UserAgent == "" {
			scenario.UserAgent, err = bot.conf.RandomUserAgent(false)
			if err != nil {
				return err
			}
		}

		// New visit
		visit, err := NewVisit(bot.searchToken, bot.analyticsToken, scenario.UserAgent, bot.conf)
		if err != nil {
			return err
		}

		// Setup specific stuff for NTO
		//visit.SetupNTO()
		// Use this line instead outside of NTO
		visit.SetupGeneral()
		visit.LastQuery.CQ = bot.conf.GlobalFilter

		err = visit.ExecuteScenario(*scenario, bot.conf)
		if err != nil {
			return err
		}

		visit.UAClient.DeleteVisit()
		time.Sleep(time.Duration(bot.random.Intn(timeVisits)) * time.Second)

		count++
		Info.Printf("Scenarios executed : %d \n =============================\n\n", count)
	}
}

func refreshScenarios(url string, isUrl bool) *Config {
	Info.Println("Updating Scenario file")

	var err error
	var conf *Config

	if isUrl {
		conf, err = NewConfigFromURL(url)
	} else {
		conf, err = NewConfigFromPath(url)
	}

	if err != nil {
		Warning.Println("Cannot update scenario file, keeping the old one")
		return nil
	} else {
		return conf
	}
}
