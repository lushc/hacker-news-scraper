package main

import (
	"context"
	"embed"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/lushc/hacker-news-scraper/internal/presenter"
	pb "github.com/lushc/hacker-news-scraper/protobufs"
)

var (
	//go:embed public/*
	content embed.FS
)

func main() {
	client, err := presenter.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()
	allBroker := presenter.NewBroker(client.WrapListAll(ctx, &emptypb.Empty{}))
	jobBroker := presenter.NewBroker(client.WrapListType(ctx, &pb.TypeRequest{Type: pb.Type_JOB}))
	storyBroker := presenter.NewBroker(client.WrapListType(ctx, &pb.TypeRequest{Type: pb.Type_STORY}))

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// SSE endpoints
	e.GET("/events/all", allBroker.Handle(ctx))
	e.GET("/events/jobs", jobBroker.Handle(ctx))
	e.GET("/events/stories", storyBroker.Handle(ctx))

	// serve embedded static content
	contentHandler := echo.WrapHandler(http.FileServer(http.FS(content)))
	contentRewrite := middleware.Rewrite(map[string]string{"/*": "/public/$1"})
	e.GET("/*", contentHandler, contentRewrite)

	// TODO: graceful shutdown
	e.Logger.Fatal(e.Start(":80"))
}
