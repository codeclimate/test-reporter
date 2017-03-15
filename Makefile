PANDOC = $(shell which pandoc)
MAN_FILES = $(wildcard man/*.md)
MAN_PAGES = $(patsubst man/%.md,man/%,$(MAN_FILES))

VERSION = 1.0.0
BUILD_VERSION = `git log -1 --pretty=format:'%H'`
BUILD_TIME = `date +%FT%T%z`
LDFLAGS = -ldflags "-X github.com/codeclimate/test-reporter/cmd.Version=${VERSION} -X github.com/codeclimate/test-reporter/cmd.BuildVersion=${BUILD_VERSION} -X github.com/codeclimate/test-reporter/cmd.BuildTime=${BUILD_TIME}"

PROJECT = /src/github.com/codeclimate/test-reporter

man/%: man/%.md
	$(PANDOC) -s -t man $< -o $@

all: $(MAN_PAGES)

test:
	go test `go list ./... | grep -v /vendor/`

build:
	go build -v ${LDFLAGS} -o $(PREFIX)bin/test-reporter$(BINARY_SUFFIX)

build-all:
	$(MAKE) build-docker GOOS=darwin GOARCH=amd64
	$(MAKE) build-docker GOOS=linux GOARCH=amd64

test-docker:
	docker run --rm \
	  --env GOPATH=/ \
	  --volume "$(PWD)":"$(PROJECT)":ro \
	  --workdir "$(PROJECT)" \
	  golang:1.8 make test

build-docker:
	docker run --rm \
	  --env PREFIX=/artifacts/ \
	  --env BINARY_SUFFIX=-$(VERSION)-$$GOOS-$$GOARCH \
	  --env GOARCH \
	  --env GOOS \
	  --env GOPATH=/ \
	  --volume "$(PWD)"/artifacts:/artifacts \
	  --volume "$(PWD)":"$(PROJECT)":ro \
	  --workdir "$(PROJECT)" \
	  golang:1.8 make build

test-ruby:
	docker build -f examples/ruby/Dockerfile .
