ifndef GOROOT
$(error GOROOT is not set)
endif

export GOPATH=$(shell pwd)

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

GO=$(GOROOT)/bin/go
GOBUILD=$(GO) build
GOINSTALL=$(GO) install

ORG_ZSTACK=org.zstack
LIBRARIES=$(ORG_ZSTACK)/server
TARGET_DIR=target


build: $(SOURCES)
	for lib in $(LIBRARIES); do \
		$(GOINSTALL) $$lib; \
	done 

.PHONY: install clean

clean:
	rm -rf target/

SERVER_BUILD_DIR=target/build/server
build-server:
	mkdir -p $(SERVER_BUILD_DIR)
	$(GOBUILD) -o $(SERVER_BUILD_DIR)/zanker-server $(ORG_ZSTACK)/server/main


install: build-server
