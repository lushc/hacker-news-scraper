.PHONY: start stop clean build consume

start:
	docker-compose up --scale consumer=0

stop:
	docker-compose stop

clean:
	docker-compose down -v --rmi all

build:
	docker-compose build

consume:
	docker-compose run consumer