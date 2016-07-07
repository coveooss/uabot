package server

import (
	"encoding/json"
	"fmt"
	"github.com/adambbolduc/uabot/autobot"
	"github.com/adambbolduc/uabot/explorerlib"
	"github.com/erocheleau/uabot/scenariolib"
	"github.com/goinggo/workpool"
	"github.com/gorilla/mux"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

var (
	quitChannels []chan bool
	workPool     *workpool.WorkPool
	random       *rand.Rand
)

type StatusResponse struct {
	Stat string `json:"status"`
}

func Status(writter http.ResponseWriter, request *http.Request) {
	status := StatusResponse{Stat: "UP"}
	json.NewEncoder(writter).Encode(status)
}

func Start(writter http.ResponseWriter, request *http.Request) {
	config := &explorerlib.Config{}
	err := json.NewDecoder(request.Body).Decode(config)
	if err != nil {
		scenariolib.Error.Println(err)
	}
	scenariolib.Info.Printf("Config : %v\n", config)
	worker := BotWorker{
		bot: autobot.NewAutobot(config, random),
	}
	err = workPool.PostWork("routine", &worker)
	if err != nil {
		fmt.Printf("Error : %v\n", err)
	}
	json.NewEncoder(writter).Encode(quitChannels)
}

func Stop(writter http.ResponseWriter, request *http.Request) {
	Vars := mux.Vars(request)
	id, _ := strconv.Atoi(Vars["id"])
	quitChannels[id] <- true
	quitChannels[id] = nil
}

func Init(_workPool *workpool.WorkPool, _random *rand.Rand) {
	workPool = _workPool
	quitChannels = make([]chan bool, 8)
	scenariolib.InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	random = _random
}

type BotWorker struct {
	bot *autobot.Autobot
}

func (worker *BotWorker) DoWork(goRoutine int) {
	quitChannel := make(chan bool)
	quitChannels[goRoutine] = quitChannel
	scenariolib.Info.Printf("Bot starting on worker: %v\n", goRoutine)
	err := worker.bot.Run(quitChannel)
	if err != nil {
		scenariolib.Error.Println(err)
	}
}
