.PHONY: build
build:
	go build -v ./cmd/todoapp

.DEFAULT_GOAL := build