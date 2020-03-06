
SHELL := /bin/bash

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

#SRC = "."

DOCKERREPO ?= $(TARGET)
DOCKERTAG ?= latest

.PHONY: all build clean install uninstall fmt simplify check run lint vet

all: check install

$(TARGET): $(SRC)
	@go build $(LDFLAGS) -o $(TARGET) ./main.go

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

docker:
	@GOOS=$(TARGETOS) make build
	@mv $(TARGET) ./dockerfile
	@docker build -t $(DOCKERREPO):$(DOCKERTAG) ./dockerfile/
	@rm ./dockerfile/$(TARGET)
	@echo New Docker image created

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
