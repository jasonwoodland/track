PHONY: build

build:
	go build

run:
	go run .

install:
	GOBIN=/usr/local/bin go install .
