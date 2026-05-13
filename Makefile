BINARY := tries
PREFIX ?= /usr/local
BINDIR := $(PREFIX)/bin
GOBIN ?= $(shell go env GOPATH)/bin

.PHONY: all build install install-local test vet fmt clean

all: build

build:
	go build -o $(BINARY) .

install: build
	install -d $(BINDIR)
	install -m 0755 $(BINARY) $(BINDIR)/$(BINARY)

install-local:
	GOBIN=$(GOBIN) go install .

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w main.go main_test.go

clean:
	rm -f $(BINARY)
