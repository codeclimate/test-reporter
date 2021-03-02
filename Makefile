.PHONY: test-docker build-docker build-linux-cgo release test-excoveralls

AWS ?= $(shell which aws)
DOCKER_RUN ?= $(shell which docker) run --rm
GIT_PUSH ?= $(shell which git) push
GIT_TAG ?= $(shell which git) tag --sign
PANDOC ?= $(shell which pandoc)

MAN_FILES = $(wildcard man/*.md)
MAN_PAGES = $(patsubst man/%.md,man/%,$(MAN_FILES))

PROJECT = github.com/codeclimate/test-reporter
VERSION ?= 0.10.0
BUILD_VERSION = $(shell git log -1 --pretty=format:'%H')
BUILD_TIME = $(shell date +%FT%T%z)
LDFLAGS = -ldflags "-X $(PROJECT)/version.Version=${VERSION} -X $(PROJECT)/version.BuildVersion=${BUILD_VERSION} -X $(PROJECT)/version.BuildTime=${BUILD_TIME}"

man/%: man/%.md
	$(PANDOC) -s -t man $< -o $@

all: test-docker build-all $(MAN_PAGES)

test:
	go test $(shell go list ./... | grep -v /vendor/)

benchmark:
	go test -bench . $(shell go list ./... | grep -v /vendor/)

build:
	if [ -z "${BUILD_TAGS}" ]; then \
		go build -v ${LDFLAGS} -o $(PREFIX)bin/test-reporter$(BINARY_SUFFIX); \
	else \
		go build -v ${LDFLAGS} -tags ${BUILD_TAGS} -o $(PREFIX)bin/test-reporter$(BINARY_SUFFIX); \
	fi


build-linux:
	$(MAKE) build \
	  PREFIX=artifacts/ \
	  BINARY_SUFFIX=-$(VERSION)-linux-amd64 \
	  CGO_ENABLED=0

build-linux-cgo:
	$(MAKE) build \
	  PREFIX=artifacts/ \
	  BINARY_SUFFIX=-$(VERSION)-netcgo-linux-amd64 \
	  CGO_ENABLED=1 \
	  BUILD_TAGS="netcgo"

build-linux-all:
	$(MAKE) build-linux
	$(MAKE) build-linux-cgo

build-darwin:
	$(MAKE) build \
	  PREFIX=artifacts/ \
	  BINARY_SUFFIX=-$(VERSION)-darwin-amd64

build-docker: BINARY_SUFFIX ?= -$(VERSION)-$$GOOS-$$GOARCH
build-docker:
	$(DOCKER_RUN) \
	  --env PREFIX=/artifacts/ \
	  --env BINARY_SUFFIX=${BINARY_SUFFIX} \
	  --env GOARCH \
	  --env GOOS \
	  --env GOPATH=/ \
	  --env CGO_ENABLED \
	  --volume "$(PWD)"/artifacts:/artifacts \
	  --volume "$(PWD)":"/src/$(PROJECT)":ro \
	  --workdir "/src/$(PROJECT)" \
	  golang:1.15 make build BUILD_TAGS=${BUILD_TAGS}

build-docker-linux:
	$(MAKE) build-docker GOOS=linux GOARCH=amd64 CGO_ENABLED=0

build-docker-linux-cgo:
	$(MAKE) build-docker GOOS=linux GOARCH=amd64 CGO_ENABLED=1 \
		BUILD_TAGS="netcgo" BINARY_SUFFIX=-$(VERSION)-netcgo-linux-amd64

test-docker:
	$(DOCKER_RUN) \
	  --env GOPATH=/ \
	  --volume "$(PWD)":"/src/$(PROJECT)":ro \
	  --workdir "/src/$(PROJECT)" \
	  golang:1.8 make test

benchmark-docker:
	$(DOCKER_RUN) \
	  --env GOPATH=/ \
	  --volume "$(PWD)":"/src/$(PROJECT)":ro \
	  --workdir "/src/$(PROJECT)" \
	  golang:1.8 make benchmark

test-simplecov:
	docker build -f integration-tests/simplecov/Dockerfile .

test-lcov:
	docker build -f integration-tests/lcov/Dockerfile .

test-covpy:
	docker build -f integration-tests/coverage_py/Dockerfile .

test-gcov:
	docker build -f integration-tests/gcov/Dockerfile .

test-gocov:
	docker build -f integration-tests/gocov/Dockerfile .

test-clover:
	docker build -f integration-tests/clover/Dockerfile .

test-cobertura:
	docker build -f integration-tests/cobertura/Dockerfile .

test-excoveralls:
	docker build -f integration-tests/excoveralls/Dockerfile .

publish-head:
	$(AWS) s3 cp \
	  --acl public-read \
	  --recursive \
	  --exclude "*" \
	  --include "test-reporter-head-*" \
	  artifacts/bin/ s3://codeclimate/test-reporter/

publish-latest:
	$(AWS) s3 cp \
	  --acl public-read \
	  --recursive \
	  --exclude "*" \
	  --include "test-reporter-latest-*" \
	  artifacts/bin/ s3://codeclimate/test-reporter/

publish-macos-version:
	if [ "$(shell curl https://s3.amazonaws.com/codeclimate/test-reporter/test-reporter-$(VERSION)-darwin-amd64 --output /dev/null --write-out %{http_code})" -eq 403 ]; then \
	  $(AWS) s3 cp \
	    --acl public-read \
	    --recursive \
	    --exclude "*" \
	    --include "test-reporter-$(VERSION)-*" \
	    artifacts/bin/ s3://codeclimate/test-reporter/; \
	else \
	  echo "Version $(VERSION) already published"; \
	  exit 1; \
	fi

publish-linux-version:
	if [ "$(shell curl https://s3.amazonaws.com/codeclimate/test-reporter/test-reporter-$(VERSION)-linux-amd64 --output /dev/null --write-out %{http_code})" -eq 403 ]; then \
	  $(AWS) s3 cp \
	    --acl public-read \
	    --recursive \
	    --exclude "*" \
	    --include "test-reporter-$(VERSION)-*" \
	    artifacts/bin/ s3://codeclimate/test-reporter/; \
	else \
	  echo "Version $(VERSION) already published"; \
	  exit 1; \
	fi

tag:
	$(GIT_TAG) --message v$(VERSION) v$(VERSION)
	$(GIT_PUSH) origin refs/tags/v$(VERSION)

clean:
	sudo $(RM) -r ./artifacts
	$(RM) $(MAN_PAGES)


# Must be run in a OS X machine. OS X binary is build natively.
manual-release:
	$(MAKE) build-docker-linux
	$(MAKE) build-docker-linux-cgo
	$(MAKE) build-darwin
	$(MAKE) build-docker-linux VERSION=latest
	$(MAKE) build-docker-linux-cgo VERSION=latest
	$(MAKE) build-darwin VERSION=latest
	$(MAKE) publish-version
	$(MAKE) publish-latest
	$(MAKE) tag
