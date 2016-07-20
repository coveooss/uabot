package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adambbolduc/uabot/explorerlib"
	"github.com/erocheleau/uabot/scenariolib"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"math/rand"
	"net/http"
	"time"
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
		http.Error(writter, err.Error(), http.StatusTeapot)
		return
	}

	err = validateConfig(config)
	if err != nil {
		http.Error(writter, err.Error(), http.StatusBadRequest)
		return
	}

	id := uuid.NewV4()

	timer := time.NewTimer(time.Duration(config.TimeToLive) * time.Minute)
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

func validateConfig(config *explorerlib.Config) error {
	if config.OriginLevels == nil {
		return errors.New("Origin Level 1 Missing")
	} else {
		for originLevel1, originLevel2 := range config.OriginLevels {
			if len(originLevel2) == 0 {
				return errors.New("Origin Level 2 Missing for originLevel1: " + originLevel1)
			}
		}
	}
	if config.SearchEndpoint == "" {
		return errors.New("searchEndpoint Missing")
	}
	if config.SearchToken == "" {
		return errors.New("searchToken Missing")
	}
	if config.AnalyticsEndpoint == "" {
		return errors.New("analyticsEndpoint Missing")
	}
	if config.AnalyticsToken == "" {
		return errors.New("analyticsToken Missing")
	}
	if config.TimeToLive == 0 {
		return errors.New("timeToLive Missing")
	}
	if config.AverageNumberOfWordsPerQuery == 0 {
		config.AverageNumberOfWordsPerQuery = 1
	}
	if config.DocumentsExplorationPercentage == 0 {
		config.DocumentsExplorationPercentage = 0.01
	}
	if config.NumberOfQueryByLanguage == 0 {
		config.NumberOfQueryByLanguage = 10
	}
	if config.FetchNumberOfResults == 0 {
		config.FetchNumberOfResults = 1000
	}
	if config.FieldsToExploreEqually == nil || len(config.FieldsToExploreEqually) == 0 {
		config.FieldsToExploreEqually = []string{"@syssource"}
	}
	if config.OutputFilePath == "" {
		config.OutputFilePath = uuid.NewV4().String() + ".json"
	}
	if config.Org == "" {

	}
	return nil
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
