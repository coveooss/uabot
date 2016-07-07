package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"net/http"
	"github.com/adambbolduc/uabot/scenariolib"
)

var (
	channels map[uuid.UUID]chan int
)

type StatusResponse struct {
	Stat string `json:"status"`
}

func Status(writter http.ResponseWriter, request *http.Request) {
	status := StatusResponse{Stat: "UP"}
	json.NewEncoder(writter).Encode(status)
}

func Start(writter http.ResponseWriter, request *http.Request) {
	if channels == nil {
		channels = make(map[uuid.UUID]chan int)
	}

	channel := make(chan int)



	id := uuid.NewV4()
	channels[id] = channel
	json.NewEncoder(writter).Encode(id)
}

func Stop(writter http.ResponseWriter, request *http.Request) {
	Vars := mux.Vars(request)
	str_id := Vars["id"]
	id, _ := uuid.FromString(str_id)
	channels[id] <- 0
	delete(channels, id)
}
