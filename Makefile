.PHONY: build

build:
	go build -o ./build/main ./cmd/shortener/main.go
run:
	./build/main

.DEFAULT_GOAL := build