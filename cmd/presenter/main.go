package main

import (
	"context"
	"embed"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/lushc/hacker-news-scraper/internal/datastore"
	"github.com/lushc/hacker-news-scraper/internal/presenter"
)

var (
	//go:embed public/*
	content embed.FS
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	src := func() ([]datastore.Item, error) {
		// TODO: gRPC client impl
		return []datastore.Item{
			{ID: 1, Type: "job", Title: "test1", Content: "", URL: "", Score: 10, CreatedBy: "foo", CreatedAt: time.Now()},
			{ID: 2, Type: "story", Title: "test2", Content: "", URL: "", Score: 20, CreatedBy: "bar", CreatedAt: time.Now()},
			{ID: 3, Type: "story", Title: "test3", Content: "", URL: "", Score: 30, CreatedBy: "baz", CreatedAt: time.Now()},
		}, nil
	}

	// SSE endpoints
	e.GET("/events/all", presenter.NewBroker(src).Handle(context.TODO()))

	// serve embedded static content
	contentHandler := echo.WrapHandler(http.FileServer(http.FS(content)))
	contentRewrite := middleware.Rewrite(map[string]string{"/*": "/public/$1"})
	e.GET("/*", contentHandler, contentRewrite)

	// TODO: graceful shutdown
	e.Logger.Fatal(e.Start(":80"))
}
