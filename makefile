MAKEFLAGS += --warn-undefined-variables
SHELL := /bin/bash
.DEFAULT_GOAL := build

.PHONY: clean lint

ROOT := $(shell pwd)
PACKAGE := HAWAI/repos/hawai-logginghub
PACKAGE_NAME := hawai-logginghub

clean:
	rm -rf build cover
	rm $(PACKAGE_NAME)

build:
	go build -v

rebuild: clean build

test:
	go test -v -race ./...

lint:
	go vet ./...
	golint ./...

all: clean build lint test
