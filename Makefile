.PHONY: test-docker build-docker build-all build-all-latest release

AWS ?= $(shell which aws)
DOCKER_RUN ?= $(shell which docker) run --rm
GIT_PUSH ?= $(shell which git) push
GIT_TAG ?= $(shell which git) tag --sign
PANDOC ?= $(shell which pandoc)

MAN_FILES = $(wildcard man/*.md)
MAN_PAGES = $(patsubst man/%.md,man/%,$(MAN_FILES))

PROJECT = github.com/codeclimate/test-reporter
VERSION ?= 0.1.4
BUILD_VERSION = $(shell git log -1 --pretty=format:'%H')
BUILD_TIME = $(shell date +%FT%T%z)
LDFLAGS = -ldflags "-X $(PROJECT)/version.Version=${VERSION} -X $(PROJECT)/version.BuildVersion=${BUILD_VERSION} -X $(PROJECT)/version.BuildTime=${BUILD_TIME}"

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

build-all-latest:
	$(MAKE) build-all VERSION=latest

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

test-simplecov:
	docker build -f examples/simplecov/Dockerfile .

test-lcov:
	docker build -f examples/lcov/Dockerfile .

test-covpy:
	docker build -f examples/coverage_py/Dockerfile .

test-gocov:
	docker build -f examples/gocov/Dockerfile .

test-clover:
	docker build -f examples/clover/Dockerfile .

publish:
	$(AWS) s3 cp --acl public-read --recursive artifacts/bin/ s3://codeclimate/test-reporter/ --exclude "*" --include "test-reporter-*"

tag:
	$(GIT_TAG) --message v$(VERSION) v$(VERSION)
	$(GIT_PUSH) origin refs/tags/v$(VERSION)

clean:
	sudo $(RM) -r ./artifacts
	$(RM) $(MAN_PAGES)

release: build-all build-all-latest publish tag
