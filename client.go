package main

import (
	"github.com/adambbolduc/uabot/server"
	"log"
	"net/http"
)

func main() {
	router := server.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

