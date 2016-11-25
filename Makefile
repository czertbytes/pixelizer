# Project
PROJECT?=pixelizer
ORGANIZATION?=czertbytes
REPOSITORY?=github.com
DOCKER_SOURCE?=Docker

# Build
BUILDER_IMAGE?=czertbytes/golang-builder:latest
GO_BUILD_PARAMS?=
GO_BUILD_CMD?=go build $(GO_BUILD_PARAMS) -o bin/pixelizer cmd/pixelizer/main.go
GO_BUILD_SERVER_CMD?=go build $(GO_BUILD_PARAMS) -o bin/server cmd/server/main.go

# Build Linux
GO_BUILD_LINUX_PARAMS?=
GO_BUILD_SERVER_LINUX_CMD?=go build $(GO_BUILD_LINUX_PARAMS) -o $(DOCKER_SOURCE)/bin/server-linux cmd/server/main.go

all: image

clean:
	rm -rf bin/
	rm -rf $(DOCKER_SOURCE)/bin

## Pixelize CLI tool
build-cli: clean
	mkdir -p bin
	$(GO_BUILD_CMD)

## Pixelize Server
build-server: clean
	mkdir -p $(DOCKER_SOURCE)/bin
	$(GO_BUILD_SERVER_CMD)

build-server-linux: clean
	mkdir -p $(DOCKER_SOURCE)/bin
	docker run --rm \
		-v $(PWD):/go/src/$(REPOSITORY)/$(ORGANIZATION)/$(PROJECT) \
		-w /go/src/$(REPOSITORY)/$(ORGANIZATION)/$(PROJECT) \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		$(BUILDER_IMAGE) \
		$(GO_BUILD_SERVER_LINUX_CMD)

## Docker
image: build-server-linux
	docker build -t $(ORGANIZATION)/$(PROJECT):latest .

.PHONY: all