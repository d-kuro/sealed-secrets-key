.PHONY: build generate clean test
BINARY_NAME  = sealed-secrets-key
LDFLAGS      = -ldflags="-s -w -X \"github.com/d-kuro/sealed-secrets-key/cmd.Revision=$(shell git rev-parse --short HEAD)\""

export GO111MODULE=on

build:
	@go build -o ./dist/$(BINARY_NAME) -v $(LDFLAGS)
