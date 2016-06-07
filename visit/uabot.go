package visit

import (
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

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

	timeNow := time.Now()
}
