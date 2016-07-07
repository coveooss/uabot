package main

import (
	"github.com/adambbolduc/uabot/server"
	"log"
	"net/http"
	"github.com/goinggo/workpool"
	"runtime"
	"time"
	"math/rand"
)

var(
	workPool *workpool.WorkPool
)

func main() {

	source := rand.NewSource(int64(time.Now().Unix()))
	random := rand.New(source)


	workPool = workpool.New(runtime.NumCPU(), 100)
	server.Init(workPool, random)
	router := server.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

