ifndef GOROOT
$(error GOROOT is not set)
endif

export GOPATH=$(shell pwd)

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

GO=$(GOROOT)/bin/go
GOBUILD=$(GO) build
GOINSTALL=$(GO) install
GOGET=$(GO) get

ORG_ZSTACK=org.zstack
LIBRARIES=$(ORG_ZSTACK)/server $(ORG_ZSTACK)/client
TARGET_DIR=target

build: $(SOURCES)
	for lib in $(LIBRARIES); do \
		$(GOINSTALL) $$lib; \
	done 

.PHONY: install clean

DEPS=github.com/Sirupsen/logrus github.com/gorilla/mux

deps:
	$(GOGET) $(DEPS)


clean:
	rm -rf target/

SERVER_BUILD_DIR=target/build/server
build-server:
	mkdir -p $(SERVER_BUILD_DIR)
	$(GOBUILD) -o $(SERVER_BUILD_DIR)/zanker-server $(ORG_ZSTACK)/server/main

CLIENT_BUILD_DIR=target/build/client
build-client:
	mkdir -p $(CLIENT_BUILD_DIR)
	$(GOBUILD) -o $(CLIENT_BUILD_DIR)/zanker $(ORG_ZSTACK)/client/main


install: build-server build-client
