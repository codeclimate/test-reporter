.PHONY: test-docker build-docker build-linux-cgo release test-excoveralls

AWS ?= $(shell which aws)
SHA_SUM ?= $(shell which shasum)
GPG ?= $(shell which gpg)
TAR ?= $(shell which tar)
DOCKER_RUN ?= $(shell which docker) run --rm
PANDOC ?= $(shell which pandoc)

MAN_FILES = $(wildcard man/*.md)
MAN_PAGES = $(patsubst man/%.md,man/%,$(MAN_FILES))

PROJECT = github.com/codeclimate/test-reporter
VERSION ?= $(shell cat VERSIONING/VERSION)
BUILD_VERSION = $(shell git log -1 --pretty=format:'%H')
BUILD_TIME = $(shell date +%FT%T%z)
LDFLAGS = -ldflags "-X $(PROJECT)/version.Version=${VERSION} -X $(PROJECT)/version.BuildVersion=${BUILD_VERSION} -X $(PROJECT)/version.BuildTime=${BUILD_TIME}"
ARTIFACTS_OUTPUT = artifacts.tar.gz

define upload_artifacts
	$(AWS) s3 cp \
	  --acl public-read \
	  --recursive \
	  --exclude "*" \
	  --include "test-reporter-$(1)-*" \
	  artifacts/bin/ s3://codeclimate/test-reporter/;
endef

define gen_signed_checksum
	cd artifacts/bin && \
	  $(SHA_SUM) -a 256 test-reporter-$(VERSION)-$(1) > test-reporter-$(VERSION)-$(1).sha256 && \
	  $(GPG) --local-user $(GPG_CODECLIMATE_FINGERPRINT) --output test-reporter-$(VERSION)-$(1).sha256.sig --detach-sig test-reporter-$(VERSION)-$(1).sha256
endef

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

build-linux-arm64:
	$(MAKE) build \
		PREFIX=artifacts/ \
		BINARY_SUFFIX=-$(VERSION)-linux-arm64 \
		CGO_ENABLED=0 \
		GOARCH=arm64

build-linux-cgo:
	$(MAKE) build \
	  PREFIX=artifacts/ \
	  BINARY_SUFFIX=-$(VERSION)-netcgo-linux-amd64 \
	  CGO_ENABLED=1 \
	  BUILD_TAGS="netcgo"

build-linux-all:
	$(MAKE) build-linux
	$(MAKE) build-linux-arm64
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

build-docker-linux-arm64:
	$(MAKE) build-docker GOOS=linux GOARCH=arm64 CGO_ENABLED=0

build-docker-linux-cgo:
	$(MAKE) build-docker GOOS=linux GOARCH=amd64 CGO_ENABLED=1 \
		BUILD_TAGS="netcgo" BINARY_SUFFIX=-$(VERSION)-netcgo-linux-amd64

build-docker-windows:
	$(MAKE) build-docker GOOS=windows GOARCH=amd64 CGO_ENABLED=0

test-docker:
	$(DOCKER_RUN) \
	  --env GOPATH=/ \
	  --volume "$(PWD)":"/src/$(PROJECT)":ro \
	  --workdir "/src/$(PROJECT)" \
	  golang:1.15 make test

benchmark-docker:
	$(DOCKER_RUN) \
	  --env GOPATH=/ \
	  --volume "$(PWD)":"/src/$(PROJECT)":ro \
	  --workdir "/src/$(PROJECT)" \
	  golang:1.15 make benchmark

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

test-dotcover:
	docker build -f integration-tests/dotcover/Dockerfile .

publish-head:
	$(call upload_artifacts,head)

publish-latest:
	$(call upload_artifacts,latest)

publish-version:
	$(call upload_artifacts,$(VERSION))

gen-linux-checksum:
	$(call gen_signed_checksum,linux-amd64)

gen-linux-arm64-checksum:
	$(call gen_signed_checksum,linux-arm64)

gen-linux-cgo-checksum:
	$(call gen_signed_checksum,netcgo-linux-amd64)

gen-darwin-checksum:
	$(call gen_signed_checksum,darwin-amd64)

gen-windows-checksum:
	$(call gen_signed_checksum,windows-amd64)

clean:
	sudo $(RM) -r ./artifacts
	$(RM) $(MAN_PAGES)

tag:
	$(TAR) -c -f ${ARTIFACTS_OUTPUT} ./artifacts/bin/test-reporter-${VERSION}-* && \
	  hub release create -a ${ARTIFACTS_OUTPUT} -m "v${VERSION}" ${VERSION}

# Must be run in a OS X machine. OS X binary is build natively.
manual-release:
	$(MAKE) build-docker-linux
	$(MAKE) build-docker-linux-arm64
	$(MAKE) build-docker-linux-cgo
	$(MAKE) build-docker-windows
	$(MAKE) build-darwin
	$(MAKE) gen-linux-checksum
	$(MAKE) gen-linux-arm64-checksum
	$(MAKE) gen-linux-cgo-checksum
	$(MAKE) gen-windows-checksum
	$(MAKE) gen-darwin-checksum
	$(MAKE) build-docker-linux VERSION=latest
	$(MAKE) build-docker-linux-arm64 VERSION=latest
	$(MAKE) build-docker-linux-cgo VERSION=latest
	$(MAKE) build-docker-windows VERSION=latest
	$(MAKE) build-darwin VERSION=latest
	$(MAKE) gen-linux-checksum VERSION=latest
	$(MAKE) gen-linux-arm64-checksum VERSION=latest
	$(MAKE) gen-linux-cgo-checksum VERSION=latest
	$(MAKE) gen-windows-checksum VERSION=latest
	$(MAKE) gen-darwin-checksum VERSION=latest
	$(MAKE) publish-version
	$(MAKE) publish-latest
