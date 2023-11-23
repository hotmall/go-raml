PACKAGES = $(shell go list ./... | grep -v '/vendor/')

PACKAGE = github.com/hotmall/go-raml
COMMIT_HASH = $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE = $(shell date +%FT%T%z)

ldflagsversion = -X $(PACKAGE)/commands.CommitHash=$(COMMIT_HASH) -X $(PACKAGE)/commands.BuildDate=$(BUILD_DATE) -s -w

all: install

install:
	go generate
	go install -v -ldflags '$(ldflagsversion)'

test:
	go test $(PACKAGES)
gogentest:
	cd codegen/golang/gentest; bash test.sh

pygentest:
	cd codegen/python/gentest; bash test.sh
