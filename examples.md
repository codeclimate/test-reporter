## The test reporter can be configured with several different languages, coverage tools, and CI's. This file contains several different working configurations.

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


- Language: Python
- CI: TravisCI
- Coverage Tool: 
- File: travis.yml
- Single/Parallel: 
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
    - CODECLIMATE_REPO_TOKEN=[token]
    - CC_TEST_REPORTER_ID=[id]

# command to run tests
```

- Language: Python
- CI: TravisCI
- Coverage Tool: Codecov
- File: travis.yml
- Single/Parallel: 
- OSS Repo: https://github.com/ukBaz/python-bluezero

```
language: python
sudo: required

# look at https://github.com/pypa/pip for examples?
# https://docs.travis-ci.com/user/multi-os/
python:
  - "3.6"
# virtualenv:
  # system_site_packages: true
# matrix:
  # include:
    # - python: 2.7
    #   env: PYTHONPATH="$PYTHONPATH:/usr/lib/python2.7/dist-packages"
    # - python: 3.3
    #   env: PYTHONPATH="$PYTHONPATH:/usr/lib/python3/dist-packages"
# addons:
#   apt:
#     packages:
#       - python-dbus
#       - python3-dbus
#       - python-gi
#       - python3-gi
before_install:
  # sudo apt-get update -qq
  # sudo apt-get install -qq python-dbus python3-dbus python-gi python3-gi
  # install dbusmock from github
  # ./install_dbusmock.sh
install:
  # - pip install --upgrade pip
  - pip install pycodestyle
  - pip install 'coverage>4.0,<4.4'
  - pip install codecov
  # Install released version of dbusmock
  # - pip install python-dbusmock
before_script:
  # - .travis/setup.sh
  # - echo "Travis Python Version ***********"
  # - echo $TRAVIS_PYTHON_VERSION
  # - echo $PYTHONPATH
  # If dbusmock installed from github
  # - export PYTHONPATH=$PYTHONPATH:/tmp/python-dbusmock-bluez_gatt
script:
  # Shared
  - coverage run -m unittest -v tests.test_tools
  - coverage run --append -m unittest -v tests.test_dbus_tools
  - coverage run --append -m unittest -v tests.test_async_tools
  # Level 100
  - coverage run --append -m unittest -v tests.test_adapter
  - coverage run --append -m unittest -v tests.test_advertisement
  - coverage run --append -m unittest -v tests.test_device
  - coverage run --append -m unittest -v tests.test_gatt
  # Level 10
  - coverage run --append -m unittest -v tests.test_broadcaster
  - coverage run --append -m unittest -v tests.test_central
  - coverage run --append -m unittest -v tests.test_peripheral
  # Level 1
  - coverage run --append -m unittest -v tests.test_eddystone
  - coverage run --append -m unittest -v tests.test_microbit
  # Examples (Level 1)
  - coverage run --append -m unittest -v tests.test_adapter_example


  - "pycodestyle bluezero"
  - "pycodestyle examples"
  # - "pycodestyle tests"
after_success:
  #
  - codecov
```
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

- Language: PHP
- CI: TravisCI
- Coverage Tool: Clover
- File: travis.yml
- Single/Parallel: Single
- OSS Repo: https://github.com/elephantly/AmpConverterBundle

```
dist: trusty
language: php
sudo: required

notifications:
    email: false

php:
  - '5.4'
  - '5.5'
  - '5.6'
  - '7.0'
  - '7.1'
  - hhvm
  - nightly

addons:
    code_climate:
        repo_token: f247b0d7edc57539b62b4be98713f14f858c80b62df7e894db65e1b4fe705027

before_install:
    - composer install

before_script:
  - curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter
  - chmod +x ./cc-test-reporter

script:
    - bin/kahlan

jobs:
  include:
    - stage: Code Climate Coverage
      env:
        global:
            - GIT_COMMITTED_AT=$(if [ "$TRAVIS_PULL_REQUEST" == "false" ]; then git log -1 --pretty=format:%ct; else git log -1 --skip 1 --pretty=format:%ct; fi)
      script:
        - ./cc-test-reporter before-build
        - bin/kahlan --config=kahlan-config.travis.php --clover=clover.xml
        - cat clover.xml
        - ./cc-test-reporter after-build --debug --exit-code $TRAVIS_TEST_RESULT
    - stage: Codecov Coverage
      script:
        - bin/kahlan --config=kahlan-config.travis.php --clover=clover.xml
        - bash <(curl -s https://codecov.io/bash)
```

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
- Language: Ruby
- CI: TravisCI
- Coverage Tool: 
- File: travis.yml
- Single/Parallel: 
- OSS Repo: https://github.com/noahd1/brakeman

```
before_install:
  - gem update bundler

script:
  - "bundle exec ruby test/test.rb"
  - "bundle exec codeclimate-test-reporter"

branches:
  only:
    - master

rvm:
  - "1.9.3"
  - "2.2.7"
  - "2.3.4"
  - "2.4.1"

addons:
  code_climate:
    repo_token: 521d341f3320acda1902d0db0a3a92fb16b11ebfe3d5ab730218d4fc0fb3db13

sudo: false
```

- Language: Ruby
- CI: TravisCI
- Coverage Tool: RSpec
- File: travis.yml
- Single/Parallel: Parallel
- OSS Repo: https://github.com/cloudfoundry/cloud_controller_ng

```
# The secure URLs below are generated using the following command:
#
# $> gem install travis
# $> travis -v
# $> travis login
# $> travis encrypt --org ENV_VAR_TO_ENCRYPT_NAME=env_var_to_encrypt_value -r cloudfoundry/cloud_controller_ng

language: ruby
bundler_args: --deployment --without development
cache: bundler
sudo: required

rvm:
  - 2.4.2

before_install:
  - gem update --system
  - wget https://github.com/nats-io/gnatsd/releases/download/v0.9.4/gnatsd-v0.9.4-linux-amd64.zip -O /tmp/gnatsd.zip
  - unzip /tmp/gnatsd.zip
  - export PATH=$PATH:$PWD/gnatsd-v0.9.4-linux-amd64

before_script:
  - bundle exec rake db:create
  - DB=mysql bundle exec rake parallel:create
  - DB=postgres bundle exec rake parallel:create
  - curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter
  - chmod +x ./cc-test-reporter
  - ./cc-test-reporter before-build

script:
  - bundle exec rake $TASKS

after_script:
  - ./cc-test-reporter after-build --exit-code $TRAVIS_TEST_RESULT

services:
  - mysql
  - postgresql

env:
  global:
    - CC_TEST_REPORTER_ID=301facccb751b8f202e8a382e9f74bda51055f738691cf2ee9a9b853ac807304
    - CF_RUN_PERM_SPECS=false

  matrix:
    - COVERAGE=true DB=postgres TASKS=spec:all
    - DB=mysql TASKS=spec:all
    - TASKS=rubocop
```

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
