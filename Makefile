PHONY: build

all: build

build:
	go build cmd/track/main.go

run:
	go run .

install:
	GOBIN=/usr/local/bin go install .
