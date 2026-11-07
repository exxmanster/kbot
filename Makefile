APP      := kbot
REGISTRY ?= ghcr.io/exxmanster
VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS  := -s -w -X github.com/exxmanster/kbot/cmd.appVersion=$(VERSION)

# Binaries: the host OS/arch by default, override with TARGETOS=... TARGETARCH=...
TARGETOS   ?= $(shell go env GOOS)
TARGETARCH ?= $(shell go env GOARCH)
OUT        ?= $(APP)

# Images: the platform of the Docker host by default, e.g. linux/arm64 on an
# Apple Silicon Mac or linux/amd64 on a usual Linux box.
IMAGE_OS   ?= $(shell docker version --format '{{.Server.Os}}')
IMAGE_ARCH ?= $(shell docker version --format '{{.Server.Arch}}')
IMAGE_TAG   = $(REGISTRY)/$(APP):$(VERSION)-$(IMAGE_OS)-$(IMAGE_ARCH)

.PHONY: help format lint test build run linux arm macos windows \
        image image-linux image-arm image-windows push clean

help:
	@echo "make format        - gofmt the code"
	@echo "make lint          - go vet"
	@echo "make test          - run unit tests"
	@echo "make build         - build ./$(APP) for this machine"
	@echo "make run           - build and start the bot (needs TELE_TOKEN)"
	@echo ""
	@echo "make linux         - dist/$(APP)-linux-amd64"
	@echo "make arm           - dist/$(APP)-linux-arm64"
	@echo "make macos         - dist/$(APP)-darwin-arm64 (Apple Silicon)"
	@echo "make windows       - dist/$(APP)-windows-amd64.exe"
	@echo ""
	@echo "make image         - image for the Docker host platform: $(IMAGE_TAG)"
	@echo "make image-linux   - linux/amd64 image"
	@echo "make image-arm     - linux/arm64 image (also runs natively on Apple Silicon Macs)"
	@echo "make image-windows - windows/amd64 image (runs on a Windows Docker host)"
	@echo "make push          - push the image to $(REGISTRY)"
	@echo "make clean         - remove binaries and the image"

format:
	gofmt -s -w .

lint:
	go vet ./...

test:
	go test ./...

build:
	CGO_ENABLED=0 GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) go build -ldflags "$(LDFLAGS)" -o $(OUT) .

run: build
	./$(APP) start

linux:
	$(MAKE) build TARGETOS=linux TARGETARCH=amd64 OUT=dist/$(APP)-linux-amd64

arm:
	$(MAKE) build TARGETOS=linux TARGETARCH=arm64 OUT=dist/$(APP)-linux-arm64

macos:
	$(MAKE) build TARGETOS=darwin TARGETARCH=arm64 OUT=dist/$(APP)-darwin-arm64

windows:
	$(MAKE) build TARGETOS=windows TARGETARCH=amd64 OUT=dist/$(APP)-windows-amd64.exe

image:
	docker build --platform $(IMAGE_OS)/$(IMAGE_ARCH) --build-arg VERSION=$(VERSION) -t $(IMAGE_TAG) .

image-linux:
	$(MAKE) image IMAGE_OS=linux IMAGE_ARCH=amd64

image-arm:
	$(MAKE) image IMAGE_OS=linux IMAGE_ARCH=arm64

image-windows:
	$(MAKE) image IMAGE_OS=windows IMAGE_ARCH=amd64

push:
	docker push $(IMAGE_TAG)

clean:
	rm -rf $(APP) dist
	-docker rmi $(IMAGE_TAG)
