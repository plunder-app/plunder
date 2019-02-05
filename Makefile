
SHELL := /bin/bash

# The name of the executable (default is current directory name)
TARGET := plunder
.DEFAULT_GOAL: $(TARGET)

# These will be provided to the target
VERSION := 0.3
BUILD := `git rev-parse HEAD`

# Operating System Default (LINUX)
TARGETOS=linux

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -s"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

DOCKERTAG=latest

.PHONY: all build clean install uninstall fmt simplify check run

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

uninstall: clean
	@rm -f $$(which ${TARGET})

fmt:
	@gofmt -l -w $(SRC)

docker:
	@GOOS=$(TARGETOS) make build
	@mv $(TARGET) ./dockerfile
	@docker build -t $(TARGET):$(DOCKERTAG) ./dockerfile/
	@rm ./dockerfile/$(TARGET)
	@echo New Docker image created

plugins:
	@echo building plugins
	@go build -buildmode=plugin -o ./plugin/example.plugin ./plugin/example.go
	@go build -buildmode=plugin -o ./plugin/kubeadm.plugin ./plugin/kubeadm/*

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
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done
	@go tool vet ${SRC}

run: install
	@$(TARGET)
