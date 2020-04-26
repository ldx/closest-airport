GIT_VERSION=$(shell git describe --dirty)
CURRENT_TIME=$(shell date +%Y%m%d%H%M%S)

LD_VERSION_FLAGS=-X main.buildVersion=$(GIT_VERSION) -X main.buildTime=$(CURRENT_TIME)
LDFLAGS=-ldflags "$(LD_VERSION_FLAGS)"

BINARIES=closest-airport

TOP_DIR=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))
SRC=$(shell find $(TOP_DIR) -type f -name '*.go')
GENERATED_SRC=airports.go

all: $(BINARIES)

airports.go: airports-data.json airports.tmpl
	sed "s/__AIRPORTS_DATA__/$(cat airports-data.json)/" airports.tmpl > $@

closest-airport: $(SRC) $(GENERATED_SRC)
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(TOP_DIR)$@

clean:
	rm -f $(BINARIES) $(GENERATED_SRC)

.PHONY: all clean
