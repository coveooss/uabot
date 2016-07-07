package main

import (
	"encoding/json"
	"flag"
	"github.com/adambbolduc/uabot/explorerlib"
	"github.com/adambbolduc/uabot/autobot"
	"github.com/erocheleau/uabot/scenariolib"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

var (
	configFile = flag.String("config", "config.json", "File Path for configuration File")
)

func main() {
	scenariolib.InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	source := rand.NewSource(int64(time.Now().Unix()))
	random := rand.New(source)

	scenariolib.Info.Println("Reading config file")
	var config explorerlib.Config
	bytes, readErr := ioutil.ReadFile(*configFile)
	check(readErr)
	json.Unmarshal(bytes, &config)

	autobot := autobot.NewAutobot(&config, random)

	stop := make(chan bool)

	go autobot.Run(stop)

	time.Sleep(5 * time.Minute)

	stop <- true


}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
