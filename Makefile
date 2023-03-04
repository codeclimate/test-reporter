.PHONY: test-podman build-podman build-linux-cgo release test-excoveralls

AWS ?= $(shell which aws)
SHA_SUM ?= $(shell which shasum)
GPG ?= $(shell which gpg)
TAR ?= $(shell which tar)
PODMAN_RUN ?= $(shell which podman) run --rm
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

all: test-podman build-all $(MAN_PAGES)

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

build-podman: BINARY_SUFFIX ?= -$(VERSION)-$$GOOS-$$GOARCH
build-podman:
	$(PODMAN_RUN) \
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

build-podman-linux:
	$(MAKE) build-podman GOOS=linux GOARCH=amd64 CGO_ENABLED=0

build-podman-linux-arm64:
	$(MAKE) build-podman GOOS=linux GOARCH=arm64 CGO_ENABLED=0

build-podman-linux-cgo:
	$(MAKE) build-podman GOOS=linux GOARCH=amd64 CGO_ENABLED=1 \
		BUILD_TAGS="netcgo" BINARY_SUFFIX=-$(VERSION)-netcgo-linux-amd64

test-podman:
	$(PODMAN_RUN) \
	  --env GOPATH=/ \
	  --volume "$(PWD)":"/src/$(PROJECT)":ro \
	  --workdir "/src/$(PROJECT)" \
	  golang:1.15 make test

benchmark-podman:
	$(PODMAN_RUN) \
	  --env GOPATH=/ \
	  --volume "$(PWD)":"/src/$(PROJECT)":ro \
	  --workdir "/src/$(PROJECT)" \
	  golang:1.15 make benchmark

test-simplecov:
	podman build -f integration-tests/simplecov/Dockerfile .

test-lcov:
	podman build -f integration-tests/lcov/Dockerfile .

test-covpy:
	podman build -f integration-tests/coverage_py/Dockerfile .

test-gcov:
	podman build -f integration-tests/gcov/Dockerfile .

test-gocov:
	podman build -f integration-tests/gocov/Dockerfile .

test-clover:
	podman build -f integration-tests/clover/Dockerfile .

test-cobertura:
	podman build -f integration-tests/cobertura/Dockerfile .

test-excoveralls:
	podman build -f integration-tests/excoveralls/Dockerfile .

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

clean:
	sudo $(RM) -r ./artifacts
	$(RM) $(MAN_PAGES)

tag:
	$(TAR) -c -f ${ARTIFACTS_OUTPUT} ./artifacts/bin/test-reporter-${VERSION}-* && \
	  hub release create -a ${ARTIFACTS_OUTPUT} -m "v${VERSION}" ${VERSION}

# Must be run in a OS X machine. OS X binary is build natively.
manual-release:
	$(MAKE) build-podman-linux
	$(MAKE) build-podman-linux-arm64
	$(MAKE) build-podman-linux-cgo
	$(MAKE) build-darwin
	$(MAKE) gen-linux-checksum
	$(MAKE) gen-linux-arm64-checksum
	$(MAKE) gen-linux-cgo-checksum
	$(MAKE) gen-darwin-checksum
	$(MAKE) build-podman-linux VERSION=latest
	$(MAKE) build-podman-linux-arm64 VERSION=latest
	$(MAKE) build-podman-linux-cgo VERSION=latest
	$(MAKE) build-darwin VERSION=latest
	$(MAKE) gen-linux-checksum VERSION=latest
	$(MAKE) gen-linux-arm64-checksum VERSION=latest
	$(MAKE) gen-linux-cgo-checksum VERSION=latest
	$(MAKE) gen-darwin-checksum VERSION=latest
	$(MAKE) publish-version
	$(MAKE) publish-latest
