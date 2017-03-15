.PHONY: test-docker build-docker build-all

AWS ?= $(shell which aws)
DOCKER_RUN ?= $(shell which docker) run --rm
PANDOC ?= $(shell which pandoc)

MAN_FILES = $(wildcard man/*.md)
MAN_PAGES = $(patsubst man/%.md,man/%,$(MAN_FILES))

PROJECT = github.com/codeclimate/test-reporter
VERSION ?= 0.1.0-rc
BUILD_VERSION = $(shell git log -1 --pretty=format:'%H')
BUILD_TIME = $(shell date +%FT%T%z)
LDFLAGS = -ldflags "-X $(PROJECT)/cmd.Version=${VERSION} -X $(PROJECT)/cmd.BuildVersion=${BUILD_VERSION} -X $(PROJECT)/cmd.BuildTime=${BUILD_TIME}"

man/%: man/%.md
	$(PANDOC) -s -t man $< -o $@

all: test-docker build-all $(MAN_PAGES)

test:
	go test $(shell go list ./... | grep -v /vendor/)

build:
	go build -v ${LDFLAGS} -o $(PREFIX)bin/test-reporter$(BINARY_SUFFIX)

build-all:
	$(MAKE) build-docker GOOS=darwin GOARCH=amd64
	$(MAKE) build-docker GOOS=linux GOARCH=amd64

test-docker:
	$(DOCKER_RUN) \
	  --env GOPATH=/ \
	  --volume "$(PWD)":"/src/$(PROJECT)":ro \
	  --workdir "/src/$(PROJECT)" \
	  golang:1.8 make test

build-docker:
	$(DOCKER_RUN) \
	  --env PREFIX=/artifacts/ \
	  --env BINARY_SUFFIX=-$(VERSION)-$$GOOS-$$GOARCH \
	  --env GOARCH \
	  --env GOOS \
	  --env GOPATH=/ \
	  --volume "$(PWD)"/artifacts:/artifacts \
	  --volume "$(PWD)":"/src/$(PROJECT)":ro \
	  --workdir "/src/$(PROJECT)" \
	  golang:1.8 make build

test-ruby:
	docker build -f examples/ruby/Dockerfile .

publish:
	$(AWS) s3 sync --acl public-read artifacts/bin s3://codeclimate/test-reporter

clean:
	sudo $(RM) -r ./artifacts
	$(RM) $(MAN_PAGES)
