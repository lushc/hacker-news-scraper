ARG go_version=1.16

# User
FROM alpine:3.13 as user
ARG uid=10001
ARG gid=10001
RUN echo "scratchuser:x:${uid}:${gid}::/home/scratchuser:/bin/sh" > /scratchpasswd

# Certs
FROM alpine:3.13 as certs
RUN apk add -U --no-cache ca-certificates

# Builder
FROM golang:${go_version}-alpine as build
ARG app
WORKDIR /code/
ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/ ./cmd/
COPY ./internal/ ./internal/

RUN GOOS=linux CGO_ENABLED=0 GOGC=off GOARCH=amd64 go build -o "./bin/${app}" "./cmd/${app}"

# Runner
FROM scratch as app
ARG app
COPY --from=user /scratchpasswd /etc/passwd
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build "/code/bin/${app}" .