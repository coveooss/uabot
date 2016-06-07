package visit

import (
	"github.com/coveo/go-coveo/analytics"
	"github.com/coveo/go-coveo/search"
)

const (
	// JSUIVERSION Change this to the version of JSUI you want to appear to be using.
	JSUIVERSION string = "0.0.0.0;0.0.0.0"
	// TIMEBETWEENACTIONS The time in seconds to wait between the different actions inside a visit
	TIMEBETWEENACTIONS int = 5

	// DEFAULTANONYMOUSTRESHOLD The default portion of users who are anonymous
	DEFAULTANONYMOUSTRESHOLD float64 = 0.5
)

// Visit    A Visit represents the execution of one scenario on a website.
// It stores data concerning the current visit including a "connection" with the
// search endpoint and the analytics endpoint.
type Visit struct {
	SearchClient search.Client
	UAClient     analytics.Client
	LastQuery    *search.Query
	LastResponse *search.Response
	OriginLevel1 string
	OriginLevel2 string
	OriginLevel3 string
	LastTab      string
	Config       *BotConfig
	User         *User
}

// NewVisit     Creates a new visit to the search page
func NewVisit(searchtoken string, uatoken string, anonymous bool, mobile bool, c *BotConfig) (*Visit, error) {

	/*if c.AllowAnonymous {
		var treshold float64
		if c.AnonymousTreshold > 0 {
			treshold = c.AnonymousTreshold
		} else {
			treshold = DEFAULTANONYMOUSTRESHOLD
		}
		if rand.Float64() <= treshold {
			anonymous = true
			Info.Printf("Anonymous visit")
		}
	}*/

	user := generateRandomUser(c, anonymous, mobile)

	Info.Printf("New Visit from: %s %s", user.Firstname, user.Lastname)
	Info.Printf("On device: %s", user.Useragent)

	// Create the http searchClient
	searchConfig := search.Config{
		Token:     searchtoken,
		UserAgent: user.Useragent,
		Endpoint:  c.SearchEndpoint,
	}
	sClient, err := search.NewClient(searchConfig)
	if err != nil {
		return nil, err
	}

	// Create the http UAClient
	uaConfig := analytics.Config{
		Token:     uatoken,
		UserAgent: user.Useragent,
		IP:        user.IP,
		Endpoint:  c.AnalyticsEndpoint,
	}
	uaClient, err := analytics.NewClient(uaConfig)
	if err != nil {
		return nil, err
	}

	return &Visit{
		Config:       c,
		User:         user,
		SearchClient: sClient,
		UAClient:     uaClient,
	}, nil
}

// ExecuteScenario Execute a specific scenario, send the config for all the
// potential random we need to do.
func (v *Visit) ExecuteScenario(scenario Scenario) error {
	Info.Printf(">>> Executing scenario named : %s", scenario.Name)

	for i := 0; i < len(scenario.Events); i++ {
		jsonEvent := scenario.Events[i]
		event, err := CreateEvent(&jsonEvent)
		if err != nil {
			return err
		}

		err = event.Execute(v)
		if err != nil {
			return err
		}

		//WaitBetweenActions()
	}
	return nil
}
