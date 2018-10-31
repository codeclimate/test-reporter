## Example 1
- Language: Ruby
- CI: TravisCI
- Coverage Tool: RSpec
- File: travis.yml
- Single/Parallel: Single
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
    - CC_TEST_REPORTER_ID=[see here: https://docs.codeclimate.com/v1.0/docs/finding-your-test-coverage-token]
    - CF_RUN_PERM_SPECS=false

  matrix:
    - COVERAGE=true DB=postgres TASKS=spec:all
    - DB=mysql TASKS=spec:all
    - TASKS=rubocop
```

## Example 2
- Language: Ruby
- CI: TravisCI
- Coverage Tool: Istanbul/SimpleCov
- File: travis.yml
- Single/Parallel: Parallel
- OSS Repo: https://github.com/scottohara/loot

```
language: ruby

# Cache gems
cache:
  bundler: true
  directories:
    - node_modules

addons:
  chrome: stable
  #firefox: latest

env:
  global:
    CC_TEST_REPORTER_ID=[see here: https://docs.codeclimate.com/v1.0/docs/finding-your-test-coverage-token]

before_install:
  - nvm install                         # Install node version from .nvmrc
  - curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter
  - chmod +x ./cc-test-reporter
  #- export DISPLAY=:99.0                # Display number for xvfb (for headless browser testing)
  #- sh -e /etc/init.d/xvfb start        # Start xvfb (for headless browser testing)

install:
  - bundle install --without production --path=${BUNDLE_PATH:-vendor/bundle}  # Install ruby gems, excluding production only gems such as unicorn (already present by default in Travis)
  - npm install                         # Install npm dependencies
  #- npm install karma-firefox-launcher codeclimate-test-reporter

# Setup the database
before_script: bundle exec rake db:create db:migrate

# Run the test suites
script:
  - bundle exec rubocop -DESP           # Backend linting
  - bundle exec rake                    # Backend specs
  #- npm test -- --browsers Firefox      # Frontend specs
  - npm test                            # Frontend linting & specs

# Pipe the coverage data to Code Climate
after_script:
  - ./cc-test-reporter format-coverage -t simplecov -o coverage/codeclimate.backend.json coverage/backend/.resultset.json # Format backend coverage
  - ./cc-test-reporter format-coverage -t lcov -o coverage/codeclimate.frontend.json coverage/frontend/lcov.info  # Format frontend coverage
  - ./cc-test-reporter sum-coverage coverage/codeclimate.*.json -p 2                  # Sum both coverage parts into coverage/codeclimate.json
  - if [[ "$TRAVIS_TEST_RESULT" == 0 ]]; then ./cc-test-reporter upload-coverage; fi  # Upload coverage/codeclimate.json
  ```

