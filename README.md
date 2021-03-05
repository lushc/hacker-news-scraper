# hacker-news-scraper

## Usage

Make a copy the env file:

`cp .env.example .env`

Start all services except the consumer:

`make start`

Once Postgres is ready, run the consumer to seed the database:

`make consume`

Remove everything once you're finished:

`make clean`