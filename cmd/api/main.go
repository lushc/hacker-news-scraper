package main

import (
	"context"
	"log"

	"github.com/lushc/hacker-news-scraper/internal/api"
	"github.com/lushc/hacker-news-scraper/internal/datastore"
)

func main() {
	ctx := context.Background()
	conn, err := datastore.NewDBConn(ctx, 5)
	if err != nil {
		log.Fatal(err)
	}

	reader, err := api.NewCachedReader(datastore.NewDBReader(*conn))
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	srv, err := api.NewServer(reader)
	if err != nil {
		log.Fatal(err)
	}

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
