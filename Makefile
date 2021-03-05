.PHONY: start stop clean build consume

start:
	docker-compose up --scale consumer=0

stop:
	docker-compose stop

clean:
	docker-compose down -v --rmi all

build:
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	protobufs/hackernews.proto
	docker-compose build

consume:
	docker-compose run consumer