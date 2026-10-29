APP     := kbot
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X github.com/exxmanster/kbot/cmd.appVersion=$(VERSION)

TARGETOS   ?= $(shell go env GOOS)
TARGETARCH ?= $(shell go env GOARCH)

.PHONY: help format lint test build run clean

help:
	@echo "make format  - gofmt the code"
	@echo "make lint    - go vet"
	@echo "make test    - run unit tests"
	@echo "make build   - build ./$(APP) for $(TARGETOS)/$(TARGETARCH)"
	@echo "make run     - build and start the bot (needs TELE_TOKEN)"
	@echo "make clean   - remove the binary"

format:
	gofmt -s -w .

lint:
	go vet ./...

test:
	go test ./...

build:
	CGO_ENABLED=0 GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) go build -ldflags "$(LDFLAGS)" -o $(APP) .

run: build
	./$(APP) start

clean:
	rm -f $(APP)
