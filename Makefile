.PHONY: run gin build test
SHELL := /bin/bash

all: run

run:
	source .env; source .env_*; go run main.go

gin:
	gin --all --immediate run main.go

build:
	rm -f elephantsql-broker
	go build -o elephantsql-broker

test:
	source .env && GOARCH=amd64 GOOS=linux go test -v ./...
