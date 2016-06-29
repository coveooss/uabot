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

func Init(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
	Trace = log.New(traceHandle, "TRACE | ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(infoHandle, "INFO | ", log.Ldate|log.Ltime)
	Warning = log.New(warningHandle, "WARNING | ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "ERROR | ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	// Init loggers
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	// Seed Random based on current time
	source := rand.NewSource(int64(time.Now().Unix()))
	random := rand.New(source)

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

	bot := scenariolib.NewUabot(local, conf, scenarioURL, searchToken, analyticsToken, random)
	err = bot.Run()
	if err != nil {
		Error.Println(err)
	}
	pp.Println("LOG >>> DONE")
}
