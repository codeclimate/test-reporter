## Example 1
- Language: Go
- CI: TravisCI
- Coverage Tool: gocov
- File: travis.yml
- Single/Parallel: Single
- OSS Repo: https://github.com/bbqtd/rg-kit

```yml
# Set the token in Travis environment settings instead defining here.
env:
  global:
    - CC_TEST_REPORTER_ID=token

language: go

# The coverprofile for multiple packages works in go 1.10
# see https://tip.golang.org/doc/go1.10#test
go:
  - master

before_script:
  - curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter
  - chmod +x ./cc-test-reporter
  - ./cc-test-reporter before-build

script:
  - go test -coverprofile c.out ./...

after_script:
  - ./cc-test-reporter after-build -t gocov --exit-code $TRAVIS_TEST_RESULT
```

## Example 2
- Language: Go (version < 1.10)
- CI: TravisCI
- Coverage Tool: gocov
- Files: codecoverage.sh, Dockerfile, travis.yml
- Single/Parallel: Single
- OSS Repo: https://github.com/nzin/dctycoon


codecoverage.sh:
```
#!/bin/sh

./cc-test-reporter before-build 
for pkg in $(go list ./... | grep -v main); do
    go test -coverprofile=$(echo $pkg | tr / -).cover $pkg
done
echo "mode: set" > c.out
grep -h -v "^mode:" ./*.cover >> c.out
rm -f *.cover

./cc-test-reporter after-build
```

Dockerfile:
```
FROM golang:1.9-alpine3.7
MAINTAINER Jordi Riera <kender.jr@gmail.com>

RUN apk add --no-cache \
    curl \
    git \
    gcc \
    cmake \
    build-base \
    libx11-dev \
    pkgconf \
    sdl2-dev \
    sdl2_ttf-dev \
    sdl2_image-dev \
    libjpeg


WORKDIR /go/src/github.com/nzin/dctycoon/
COPY . .
RUN go get -u github.com/golang/lint/golint && \
    go get -u github.com/jteeuwen/go-bindata/... && \
    go get -u github.com/stretchr/testify/assert && \
    go get github.com/axw/gocov/gocov && \
    curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter && \
    chmod +x ./cc-test-reporter

ENV CC_TEST_REPORTER_ID=token

RUN "$(go env GOPATH)/bin/go-bindata" -o global/assets.go -pkg global assets/... && \
    go get ./... && \
    go build ./...
```

travis.yml
```yaml
services:
  - docker

script:
  - docker build -t app .
  - docker run --rm app ./codecoverage.sh
```


## Example 3
- Language: Go 1.9
- CI: TravisCI
- Coverage Tool: gocov
- File: travis.yml
- Single/Parallel: Single

```
language: go
go:
  - 1.9
install:
  - go get -v github.com/codeclimate/test-reporter
  - cd $GOPATH/src/github.com/codeclimate/test-reporter && git checkout tags/v0.4.3 && go install
  - cd -
before_script:
  - test-reporter before-build
script:
  - go test -coverprofile c.out -coverpkg ./...
after_script:
  - test-reporter after-build --exit-code $TRAVIS_TEST_RESULT
env:
  global:
    - secure: [REDACTED]
    
```
