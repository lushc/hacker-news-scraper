package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/lushc/hacker-news-scraper/internal/api"
)

const (
	serverPortEnv = "SERVER_PORT"
)

var (
	errServerPortEnv = fmt.Errorf("missing env var %s", serverPortEnv)
)

func main() {
	portEnv, ok := os.LookupEnv(serverPortEnv)
	if !ok {
		log.Fatal(errServerPortEnv)
	}

	port, err := strconv.Atoi(portEnv)
	if err != nil {
		log.Fatal(err)
	}

	if err := api.StartServer(port); err != nil {
		log.Fatal(err)
	}
}
