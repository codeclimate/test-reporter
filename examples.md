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
