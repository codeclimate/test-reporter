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
