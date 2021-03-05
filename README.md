# hacker-news-scraper

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

Remove everything once you're finished:

`make clean`