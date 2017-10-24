## Example 1
- Language: PHP
- CI: TravisCI
- Coverage Tool: PHP Codeception
- File: travis.yml
- Single/Parallel: 

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
- ./cc-test-reporter after-build --coverage-input-type clover --id 12345 --exit-code $TRAVIS_TEST_RESULT
```


## Example 2
- Language: PHP
- CI: TravisCI
- Coverage Tool: Clover
- File: travis.yml
- Single/Parallel: Single
- OSS Repo: https://github.com/trogne/skeleton

```
env:
  global:
    - CC_TEST_REPORTER_ID=7200f3ac9aab067d6a3c75ddf45f1cadbfb0ee1f9ef902f4e2d005a2511c5745
    - GIT_COMMITTED_AT=$(if [ "$TRAVIS_PULL_REQUEST" == "false" ]; then git log -1 --pretty=format:%ct; else git log -1 --skip 1 --pretty=format:%ct; fi)    
language: php
php:
  - 7.0
before_script:
  - "composer require codeclimate/php-test-reporter --dev"
  - "composer install"
  - curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter
  - chmod +x ./cc-test-reporter
  - ./cc-test-reporter before-build
script:
  - "phpunit --testsuite=unit --coverage-text --coverage-clover build/logs/clover.xml"
  - if [ "$TRAVIS_PULL_REQUEST" == "false" ]; then ./cc-test-reporter after-build --exit-code $TRAVIS_TEST_RESULT; fi
  ```


## Example 3
- Language: PHP
- CI: TravisCI
- Coverage Tool: 
- File: travis.yml
- Single/Parallel: 
- OSS Repo: https://github.com/jmwri/pubg-php

```
env:
  global:
    - CC_TEST_REPORTER_ID=86a09970f02b3f841b263963099a65adf9bf9e85ea72db88a1d5f5ac003d01f2
    - GIT_COMMITTED_AT=$(if [ "$TRAVIS_PULL_REQUEST" == "false" ]; then git log -1 --pretty=format:%ct; else git log -1 --skip 1 --pretty=format:%ct; fi)
language: php
php:
  - '5.6'
  - '7.1'
before_script:
  - curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter
  - chmod +x ./cc-test-reporter
  - if [ $(phpenv version-name) = "7.1" ]; then ./cc-test-reporter before-build; fi
install:
  - composer install
after_script:
    - if [ $(phpenv version-name) = "7.1" ] && [ "$TRAVIS_PULL_REQUEST" == "false" ]; then ./cc-test-reporter after-build --exit-code $TRAVIS_TEST_RESULT; fi
  ```
  

## Example 4
- Language: PHP
- CI: CircleCI 2.0
- Coverage Tool: PHPUnit/Clover
- File: config.yml
- Single/Parallel: 
- OSS Repo: https://github.com/ejcnet/sourcebot
- Check out this blog post for more info! https://medium.com/@paulmwatson/configuring-code-coverage-for-code-climate-with-circleci-2-0-and-phpunit-3f7612683b67

```
version: 2

jobs:
  build:
    environment:
      CC_TEST_REPORTER_ID: YOUR_CODE_CLIMATE_REPORTER_ID
    docker:
      - image: notnoopci/php:7.1.5-browsers
    working_directory: ~/repo
    steps:
      - checkout
      - run: sudo pecl channel-update pecl.php.net
      - run: sudo pecl install xdebug && sudo docker-php-ext-enable xdebug
      - run: curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter
      - run: chmod +x ./cc-test-reporter
      - run: sudo mkdir -p $CIRCLE_TEST_REPORTS/phpunit
      - run: ./cc-test-reporter before-build
      - run: sudo vendor/bin/phpunit --coverage-clover clover.xml
      - run: ./cc-test-reporter after-build -t clover --exit-code $?
  ```
