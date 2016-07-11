package scenariolib

import (
	"math/rand"
	"time"
)

// DEFAULTTIMEBETWEENVISITS The time for the bot to wait between visits, between 0 and X Seconds
const DEFAULTTIMEBETWEENVISITS int = 300

// DEFAULT_STANDARD_DEVIATION_BETWEEN_VISITS The standard deviation when updating time between visits
const DEFAULT_STANDARD_DEVIATION_BETWEEN_VISITS int = 150

// WEEKEND_MODIFIER The modifier to multiply DEFAULTTIMEBETWEENVISITS during weekends
const WEEKEND_MODIFIER = 10

type Uabot interface {
	Run() error
}

type uabot struct {
	local             bool
	scenarioURL       string
	searchToken       string
	analyticsToken    string
	random            *rand.Rand
	WaitBetweenVisits bool
}

func NewUabot(local bool, scenarioUrl string, searchToken string, analyticsToken string, random *rand.Rand) *uabot {
	return &uabot{
		local,
		scenarioUrl,
		searchToken,
		analyticsToken,
		random,
		true,
	}
}

func (bot *uabot) Run(quitChannel chan bool) error {
	var (
		conf       *Config
		err        error
		timeVisits int
	)

	// Init from path instead of URL, for testing purposes
	if bot.local {
		conf, err = NewConfigFromPath(bot.scenarioURL)
	} else {
		conf, err = NewConfigFromURL(bot.scenarioURL)
	}
	if err != nil {
		return err
	}

	bot.WaitBetweenVisits = !conf.DontWaitBetweenVisits

	// Refresh the scenario files every 5 hours automatically.
	// This way, no need to stop the bot to update the possible scenarios.
	bot.continuallyRefreshScenariosEvery(5*time.Hour, conf)
	if conf.TimeBetweenVisits > 0 {
		timeVisits = conf.TimeBetweenVisits
	} else {
		timeVisits = DEFAULTTIMEBETWEENVISITS
		bot.continuallyUpdateTimeVisitsEvery(24*time.Hour, &timeVisits)
	}

	count := 0
	for { // Run forever

		select {
		default:
			scenario, err := conf.RandomScenario()
			if err != nil {
				return err
			}

			if scenario.UserAgent == "" {
				scenario.UserAgent, err = conf.RandomUserAgent(false)
				if err != nil {
					return err
				}
			}

			// New visit
			visit, err := NewVisit(bot.searchToken, bot.analyticsToken, scenario.UserAgent, scenario.Language, conf)
			if err != nil {
				return err
			}

			// Setup specific stuff for NTO
			//visit.SetupNTO()
			// Use this line instead outside of NTO
			visit.SetupGeneral()
			visit.LastQuery.CQ = conf.GlobalFilter

			err = visit.ExecuteScenario(*scenario, conf)
			if err != nil {
				return err
			}

			visit.UAClient.DeleteVisit()
			if bot.WaitBetweenVisits {
				waitTime := time.Duration(bot.random.Intn(timeVisits)) * time.Second
				time.Sleep(waitTime)
			}

			count++
			Info.Printf("Scenarios executed : %d \n =============================\n\n", count)

		case <-quitChannel:
			return nil
		}
	}
}

func (bot *uabot) continuallyUpdateTimeVisitsEvery(timeDuration time.Duration, timeVisits *int) {
	ticker := time.NewTicker(timeDuration)
	go func() {
		for _ = range ticker.C {
			var effectiveMeanTimeBetweenVisits = DEFAULTTIMEBETWEENVISITS
			if time.Now().Weekday() == time.Saturday || time.Now().Weekday() == time.Sunday {
				effectiveMeanTimeBetweenVisits = DEFAULTTIMEBETWEENVISITS * WEEKEND_MODIFIER
			}
			var randomPositiveTime int
			for randomPositiveTime = 0; randomPositiveTime <= 0; randomPositiveTime = int(float64(DEFAULT_STANDARD_DEVIATION_BETWEEN_VISITS)*bot.random.NormFloat64()+0.5) + effectiveMeanTimeBetweenVisits {
			}
			*timeVisits = randomPositiveTime
			Info.Println("Updating Time Visits to", *timeVisits)
		}
	}()
}

func (bot *uabot) continuallyRefreshScenariosEvery(timeDuration time.Duration, conf *Config) {
	ticker := time.NewTicker(timeDuration)
	go func() {
		for _ = range ticker.C {
			conf2 := refreshScenarios(bot.scenarioURL, bot.local)
			if conf2 != nil {
				Info.Println("Refreshing scenario")
				conf = conf2
			}
		}
	}()
}

func refreshScenarios(url string, isLocal bool) *Config {
	Info.Println("Updating Scenario file")

	var err error
	var conf *Config

	if isLocal {
		conf, err = NewConfigFromPath(url)

	} else {
		conf, err = NewConfigFromURL(url)
	}

	if err != nil {
		Warning.Println("Cannot update scenario file, keeping the old one")
		return nil
	} else {
		return conf
	}
}
