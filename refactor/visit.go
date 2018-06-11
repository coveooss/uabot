package refactor

import (
	"github.com/go-coveo/analytics"
	"github.com/go-coveo/search"
)

const (
	// JSUIVERSION Change this to the version of JSUI you want to appear to be using.
	JSUIVERSION string = "0.0.0.0;0.0.0.0"
)

// A Visit is when one "person" arrives on the website until he leaves.
type Visit struct {
	SearchClient search.Client
	UAClient     analytics.Client
	LastQuery    *search.Query
	LastResponse *search.Response
	User         User
	OriginLevel1 string
	OriginLevel2 string
	OriginLevel3 string
	Config       *BotConfig
}

type User struct {
	FirstName string
	LastName  string
	Email     string
	IP        string
	UserAgent string
	Language  string
	Anonymous bool
}
