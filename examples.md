## The test reporter can be configured with several different languages, coverage tools, and CI's. This file contains several different working configurations.

- Language: PHP
- CI: TravisCI
- Testing Framework: PHP Codeception
- File: travis.yml

```
language: php

php: 
- 7.1

before_script: 
- curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter 
- chmod +x ./cc-test-reporter 
- ./cc-test-reporter before-build

script:

# Running unit tests with clover coverage report 
- vendor/bin/codecept run unit --coverage --coverage-xml

after_script: 
- mv tests/_output/coverage.xml clover.xml 
- ./cc-test-reporter after-build --coverage-input-type clover --id 12345 --exit-code $TRAVIS_TEST_RESULT```


- Language: Python
- CI: TravisCI
- Testing Framework: 
- File: travis.yml
- OSS Repo: https://github.com/menntamalastofnun/skolagatt

```
dist: trusty
language: python
python:
  - "3.5"
# command to install dependencies
before_script:
  - curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter
  - chmod +x ./cc-test-reporter
install:
  - pip install -r requirements.txt
addons:
  postgresql: "9.5"
services:
  - redis-server
# for codecoverage on codeclimate.com
env:
  global:
    - GIT_COMMITTED_AT=$(if [ "$TRAVIS_PULL_REQUEST" == "false" ]; then git log -1 --pretty=format:%ct; else git log -1 --skip 1 --pretty=format:%ct; fi)
    - CODECLIMATE_REPO_TOKEN=487c91389b4ddd1ea050bef2a1c822db7d54f302aad677daa82cd7ae92adfeb6
    - CC_TEST_REPORTER_ID=487c91389b4ddd1ea050bef2a1c822db7d54f302aad677daa82cd7ae92adfeb6

# command to run tests
```
