PANDOC = $(shell which pandoc)
MAN_FILES = $(wildcard man/*.md)
MAN_PAGES = $(patsubst man/%.md,man/%,$(MAN_FILES))

VERSION = 1.0.0
BUILD_VERSION = `git log -1 --pretty=format:'%H'`
BUILD_TIME = `date +%FT%T%z`
LDFLAGS = -ldflags "-X github.com/codeclimate/test-reporter/cmd.Version=${VERSION} -X github.com/codeclimate/test-reporter/cmd.BuildVersion=${BUILD_VERSION} -X github.com/codeclimate/test-reporter/cmd.BuildTime=${BUILD_TIME}"

man/%: man/%.md
	$(PANDOC) -s -t man $< -o $@

all: $(MAN_PAGES)

test:
	go test `go list ./... | grep -v /vendor/`

build:
	go build -v ${LDFLAGS} -o bin/test-reporter
