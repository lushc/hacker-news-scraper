# hacker-news-scraper

My take on the GS on-boarding project.

⚠️ Zero tests, sad-looking frontend ⚠️

## Requirements

* Go 1.16
* Protocol buffer compiler & gRPC (see [Quick Start](https://grpc.io/docs/languages/go/quickstart/#prerequisites))
* Docker

## Usage

Make a copy the env file:

`cp .env.example .env`

Start all services except the consumer:

`make start`

Once Postgres is ready, run the consumer to seed the database:

`make consume`

Browse to http://localhost:8080/ to view the stories (use the radio buttons to switch between types).

Remove everything once you're finished:

`make clean`

## Architecture

### Consumer

Starts concurrent workers that fetch items from the [New, Top and Best Stories](https://github.com/HackerNews/API#new-top-and-best-stories)
endpoint and inserts them into a Postgres table.

### API

Starts a gRPC server which has methods for streaming story items from Postgres through a short-lived Redis-backed caching layer.

### Presenter

Starts an Echo webserver and connects to the API as a gRPC client, concurrently streaming story items to web clients through Server-Sent Events.

Note: the data is streamed once when the web client connects, this currently doesn't react new inserts made by the consumer.