package main

import (
	"github.com/adambbolduc/uabot/server"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

func main() {
	source := rand.NewSource(int64(time.Now().Unix()))
	random := rand.New(source)

	queueLength := int32(100)
	concurrentGoRoutine := runtime.NumCPU()
	workPool := server.NewWorkPool(concurrentGoRoutine, queueLength)

	server.Init(workPool, random)
	router := server.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
