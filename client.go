package main

import (
	"github.com/adambbolduc/uabot/server"
	"log"
	"net/http"
	"github.com/goinggo/workpool"
	"runtime"
)

var(
	workPool *workpool.WorkPool
)

func main() {
	workPool = workpool.New(runtime.NumCPU(), 100)
	server.Inject(workPool)
	router := server.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

