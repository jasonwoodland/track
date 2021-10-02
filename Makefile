PHONY: build

all: build

build:
	go build ./cmd/track

install:
	GOBIN=/usr/local/bin go install ./cmd/track
