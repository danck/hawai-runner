MAKEFLAGS += --warn-undefined-variables
SHELL := /bin/bash
.DEFAULT_GOAL := build

.PHONY: clean lint

ROOT := $(shell pwd)
PACKAGE := HAWAI/repos/hawai-logginghub

clean:
	rm -rf build cover
	rm hawai-logginghub

build:
	go build -v

rebuild: clean build

test:
	go test -v -race ./...

lint:
	go vet ./...
	golint ./...

package:
	sudo docker build -t hawai-logginghub $(ROOT)

all: clean build lint test package

run:
	sudo docker run --publish 20000:20000 --name testlogging --rm hawai-logginghub
