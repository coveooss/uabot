package server

import (
	"encoding/json"
	"fmt"
	"github.com/adambbolduc/uabot/autobot"
	"github.com/erocheleau/uabot/scenariolib"
	"github.com/gorilla/mux"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

var (
	quitChannels []chan bool
	random       *rand.Rand
	workPool     *WorkPool
)

func Init(_workPool *WorkPool, _random *rand.Rand) {
	workPool = _workPool
	quitChannels = make([]chan bool, 8)
	scenariolib.InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	random = _random
}

func Start(writter http.ResponseWriter, request *http.Request) {
	config, err := DecodeConfig(request.Body)
	if err != nil {
		http.Error(writter, err.Error(), 418)
		return
	}
	worker := WorkWrapper{
		realWorker: &BotWorker{
			bot: autobot.NewAutobot(config, random),
		},
		workPool: workPool,
	}
	err = workPool.PostWork(&worker)
	if err != nil {
		fmt.Printf("Error : %v\n", err)
	}
}

func Stop(writter http.ResponseWriter, request *http.Request) {
	Vars := mux.Vars(request)
	id, _ := strconv.Atoi(Vars["id"])
	quitChannels[id] <- true
	quitChannels[id] = nil
}

func GetInfo(writter http.ResponseWriter, request *http.Request) {
	infos := map[string]interface{}{
		"status":         "UP",
		"botWorkerInfos": workPool.getInfo(),
		"activeRoutines": fmt.Sprintf("%v/%v", workPool.ActiveRoutines(), workPool.NumberConcurrentRoutine),
		"queuedWork":     fmt.Sprintf("%v/%v", workPool.QueuedWork(), workPool.QueueLength),
	}
	json.NewEncoder(writter).Encode(infos)
}
