# Makefile

NAME := st-proc
SRC_URL := github.com/dtroyer/$(NAME)

BIN_DIR ?= bin
SHELL := /bin/bash

.PHONY: build
build: ${BIN_DIR}
	go build -o $(BIN_DIR)/$(NAME)

.PHONY: clean
clean:
	rm -rf ${BIN_DIR}/*

.PHONY: fmt
fmt:
	files=$$(gofmt -l . | tee >(cat - >&2)); [ -z "$$files" ]

.PHONY: fmtfix
fmtfix:
	go fmt ./...

.PHONY: setup
setup:
	go get -t $(SRC_URL)/cmd $(SRC_URL)/flight

.PHONY: test
test:
	go test ./... $(ARGS)

testc:
	$(MAKE) test ARGS=-cover

${BIN_DIR}:
	mkdir -p ${BIN_DIR}
