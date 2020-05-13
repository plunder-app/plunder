
SHELL := /bin/sh

# The name of the executable (default is current directory name)
TARGET := plunder
.DEFAULT_GOAL: $(TARGET)

# These will be provided to the target
VERSION := 0.5.0
BUILD := `git rev-parse HEAD`

# Required for the move to go modules for >v0.5.0
export GO111MODULE=on

# Operating System Default (LINUX)
TARGETOS=linux

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -s"

REPOSITORY = plndr
DOCKERREPO ?= $(TARGET)
DOCKERTAG ?= latest

.PHONY: all build clean install uninstall fmt simplify check run lint vet

all: check install

$(TARGET): $(SRC)
	@go build $(LDFLAGS) -o $(TARGET)

build: $(TARGET)
	@true

clean:
	@rm -f $(TARGET)

install:
	@echo Building and Installing project
	@go install $(LDFLAGS)

install_plugin:
	@make plugins
	@echo Installing plugins
	-mkdir ~/plugin
	-cp -pr ./plugin/*.plugin ~/plugin/

uninstall: clean
	@rm -f $$(which ${TARGET})

fmt:
	@gofmt -l -w $(SRC)

vet:
	@go vet $(SRC)

lint:
	@golint $(SRC)

# This is typically only for quick testing
dockerx86:
	@docker buildx build  --platform linux/amd64 --load -t $(REPOSITORY)/$(TARGET):$(DOCKERTAG) -f Dockerfile .
	@echo New Multi Architecture Docker image created

docker:
	@docker buildx build  --platform linux/amd64,linux/arm64,linux/arm/v7 --push -t $(REPOSITORY)/$(TARGET):$(DOCKERTAG) -f Dockerfile .
	@echo New Multi Architecture Docker image created

plugins:
	@echo Building plugins
	@GO111Module=off go build -buildmode=plugin -o ./plugin/example.plugin ./plugin/example.go
	@GO111Module=off go build -buildmode=plugin -o ./plugin/kubeadm.plugin ./plugin/kubeadm/*
	@GO111Module=off go build -buildmode=plugin -o ./plugin/docker.plugin ./plugin/docker/*

release_darwin:
	@echo Creating Darwin Build
	@GOOS=darwin make build
	@GOOS=darwin make plugins
	@zip -9 -r plunder-darwin-$(VERSION).zip ./plunder ./plugin/*.plugin
	@rm plunder
	@rm ./plugin/*.plugin

release_linux:
	@echo Creating Linux Build
	@GOOS=linux make build
	@GOOS=linux make plugins
	@zip -9 -r plunder-linux-$(VERSION).zip ./plunder ./plugin/*.plugin
	@rm plunder
	@rm ./plugin/*.plugin

simplify:
	@gofmt -s -l -w $(SRC)

check:
	@test -z $(shell gofmt -l main.go | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	make lint
	make vet

run: install
	@$(TARGET)
