GIT_VERSION=$(shell git describe --dirty)
CURRENT_TIME=$(shell date +%Y%m%d%H%M%S)

LD_VERSION_FLAGS=-X main.buildVersion=$(GIT_VERSION) -X main.buildTime=$(CURRENT_TIME)
LDFLAGS=-ldflags "$(LD_VERSION_FLAGS)"

BINARIES=closest-airport

TOP_DIR=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))
SRC=$(shell find $(TOP_DIR) -type f -name '*.go')
GENERATED_SRC=airports.go

all: $(BINARIES)

airports.go: data/airport-data.json
	go get -u github.com/go-bindata/go-bindata/...
	$(GOPATH)/bin/go-bindata -o airports.go data/

closest-airport: $(SRC) $(GENERATED_SRC)
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(TOP_DIR)$@

clean:
	rm -f $(BINARIES) $(GENERATED_SRC)

.PHONY: all clean
