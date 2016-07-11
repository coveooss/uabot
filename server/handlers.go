package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"math/rand"
	"net/http"
	"time"
	"github.com/erocheleau/uabot/scenariolib"
)

var (
	quitChannels map[uuid.UUID]chan bool
	random       *rand.Rand
	workPool     *WorkPool
)

func Init(_workPool *WorkPool, _random *rand.Rand) {
	workPool = _workPool
	quitChannels = make(map[uuid.UUID]chan bool)
	random = _random
}

func Start(writter http.ResponseWriter, request *http.Request) {
	config, err := DecodeConfig(request.Body)
	if err != nil {
		http.Error(writter, err.Error(), 418)
		return
	}

	id := uuid.NewV4()
	if config.OutputFilePath == "" {
		config.OutputFilePath = id.String() + ".json"
	}

	timer := time.NewTimer(time.Duration(config.Timeout) * time.Minute)
	quitChannel := make(chan bool)
	go func() {
		<-timer.C
		scenariolib.Info.Printf("Timer Timed Out")
		close(quitChannel)
	}()
	quitChannels[id] = quitChannel
	worker := NewWorker(config, quitChannel, random, id)
	err = workPool.PostWork(&worker)
	if err != nil {
		fmt.Printf("Error : %v\n", err)
	}
	json.NewEncoder(writter).Encode(map[string]interface{}{
		"workerID": id,
	})
}

func Stop(writter http.ResponseWriter, request *http.Request) {
	Vars := mux.Vars(request)
	id, _ := uuid.FromString(Vars["id"])
	close(quitChannels[id])
	delete(quitChannels, id)
}

func GetInfo(writter http.ResponseWriter, request *http.Request) {
	infos := map[string]interface{}{
		"status":         "UP",
		"botWorkerInfos": workPool.getInfo(),
		"activeRoutines": fmt.Sprintf("%v/%v", workPool.ActiveRoutines(), workPool.NumberConcurrentRoutine),
		"queuedWork":     fmt.Sprintf("%v/%v", workPool.QueuedWork(), workPool.QueueLength),
	}
	writter.Header().Add("Content-Type", "application/json")
	json.NewEncoder(writter).Encode(infos)
}
