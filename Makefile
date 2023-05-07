PHONY: build

all: build

build:
	go build ./cmd/track

install:
	go install ./cmd/track
