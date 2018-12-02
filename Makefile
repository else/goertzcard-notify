BINARY := goertzcard-notify
BUILD_VERSION := $(shell git describe --tags --always --dirty --abbrev)
BUILD_DATE := $(shell date +%FT%T%z)

LDFLAGS=-ldflags "-w -s -X main.Version=$(BUILD_VERSION) -X main.BuildDate=$(BUILD_DATE)"

build:
	go build -race ${LDFLAGS} -o ${BINARY}

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY}

build-docker:
	docker build -t else/goertzcard-notify .

test:
	go test -v ./...

install:
	go install ${LDFLAGS}

clean:
	rm -f ${BINARY}

.PHONY: build build-linux build-docker install test clean