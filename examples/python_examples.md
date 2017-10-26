## Example 1
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


## Example 2
- Language: Python
- CI: TravisCI
- Coverage Tool: 
- File: travis.yml
- Single/Parallel: 
- OSS Repo: https://github.com/czlee/tabbycat

```
language: python
dist: trusty
sudo: required
group: edge
python:
  - "3.4"
  - "3.5"
  - "3.6"
env:
  global:
    - CC_TEST_REPORTER_ID=token
addons:
  postgresql: "9.6"
  # chrome: stable # Re-enable for functional tests
services:
  - postgresql
install:
  - pip install -r requirements_common.txt
  - pip install coverage
  - npm install
before_script:  # code coverage tool
  - curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter
  - chmod +x ./cc-test-reporter
  - ./cc-test-reporter before-build
script:
  - flake8 tabbycat
  - coverage run tabbycat/manage.py test -v 2 --exclude-tag=selenium
after_script:
  - coverage xml
  - if [[ "$TRAVIS_PULL_REQUEST" == "false" && "$TRAVIS_PYTHON_VERSION" == "3.6" ]]; then ./cc-test-reporter after-build --exit-code $TRAVIS_TEST_RESULT; fi
# The below is used to enable selenium testing as per:
# https://docs.travis-ci.com/user/gui-and-headless-browsers/
# Currently it runs and loads the view; but doesn't seem to resolve the asserts
# Either the Chrome instance isn't running; or the static files aren't serving
# To rule out the former maybe disable tabbycat.standings.tests.test_ui.CoreStandingsTests
# And just let it test the login page (that should work without staticfiles)
# before_install:
#   # Run google chrome in headless mode
#   - google-chrome-stable --headless --disable-gpu --remote-debugging-port=9222 http://localhost &
# before_script:
#   # GUI for real browsers.
#   - export DISPLAY=:99.0
#   - sh -e /etc/init.d/xvfb start
#   - sleep 3 # give xvfb some time to start
#   - dj collectstatic
```
