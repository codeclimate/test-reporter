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

## Example 3
- Language: Ruby
- CI: CircleCI
- Coverage Tool: RSpec
- File: .circleci/config.yml
- Single/Parallel: Parallel
- OSS Repo: [This repo](https://github.com/FlorinPopaCodes/bike_index/blob/master/.circleci/config.yml), based on [this conversation](https://github.com/codeclimate/test-reporter/issues/386#issuecomment-473684517)

```
version: 2
jobs:
  build:
    working_directory: ~/bikeindex/bike_index
    parallelism: 4
    shell: /bin/bash --login
    environment:
      RAILS_ENV: test
      RACK_ENV: test
      COVERAGE: true
      TZ: /usr/share/zoneinfo/America/Chicago
    docker:
      - image: circleci/ruby:2.5.1-stretch-node
        environment:
          PGHOST: 127.0.0.1
          PGUSER: root
      - image: circleci/postgres:10.4-alpine-postgis
        environment:
          POSTGRES_USER: root
          POSTGRES_DB: bikeindex_test
      - image: redis:4.0.9
    steps:
      - checkout
      - restore_cache:
          keys:
            # This branch if available
            - v2-dep-{{ .Branch }}-
            # Default branch if not
            - v2-dep-master-
            # Any branch if there are none on the default branch - this should be unnecessary if you have your default branch configured correctly
            - v2-dep-
      - run:
          name: install dockerize
          command: wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz && sudo tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz
          environment:
            DOCKERIZE_VERSION: v0.6.1
      - run:
          name: install system libraries
          command: sudo apt-get update && sudo apt-get -y install imagemagick postgresql-client
      - run:
          name: install bundler
          command: gem install bundler
      - run:
          name: bundle gems
          command: bundle install --path=vendor/bundle --jobs=4 --retry=3
      # So that we can compile assets, since we use node & yarn
      - run:
          name: Yarn Install
          command: yarn install --cache-folder ~/.cache/yarn
      - run: bundle exec rake assets:precompile
      - run:
          name: Install Code Climate Test Reporter
          command: |
            curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter
            chmod +x ./cc-test-reporter
      - run:
          name: Wait for PostgreSQL to start
          command: dockerize -wait tcp://localhost:5432 -timeout 1m
      - save_cache:
          key: v2-dep-{{ .Branch }}-{{ epoch }}
          paths:
            - ./vendor/bundle
            - ~/.bundle
            - public/assets
            - tmp/cache/assets/sprockets
            - ~/.cache/yarn
      - run:
          name: Setup Database
          command: |
            bundle exec rake db:create db:structure:load
      - run:
          name: Run tests
          command: |
            mkdir -p test-results/rspec test-artifacts
            ./cc-test-reporter before-build
            TESTFILES=$(circleci tests glob "spec/**/*_spec.rb" | circleci tests split --split-by=timings)
            bundle exec rspec --profile 10 \
                              --color \
                              --order random \
                              --format RspecJunitFormatter \
                              --out test-results/rspec/rspec.xml \
                              --format progress \
                              -- ${TESTFILES}
      - run:
          name: Code Climate Test Coverage
          command: |
            ./cc-test-reporter format-coverage -t simplecov -o "coverage/codeclimate.$CIRCLE_NODE_INDEX.json"
      - persist_to_workspace:
          root: coverage
          paths:
            - codeclimate.*.json
      - store_test_results:
          path: test-results
      - store_artifacts:
          path: test-artifacts

  upload-coverage:
    docker:
      - image: circleci/ruby:2.5.1-stretch-node
    environment:
      CC_TEST_REPORTER_ID: c3ff91e23ea0fea718bb62dae0a8a5440dc082d5d2bb508af6b33d0babac479a
    working_directory: ~/bikeindex/bike_index

    steps:
      - attach_workspace:
          at: ~/bikeindex/bike_index
      - run:
          name: Install Code Climate Test Reporter
          command: |
            curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter
            chmod +x ./cc-test-reporter
      - run:
          command: |
            ./cc-test-reporter sum-coverage --output - codeclimate.*.json | ./cc-test-reporter upload-coverage --debug --input -

workflows:
  version: 2

  commit:
    jobs:
      - build
      - upload-coverage:
          requires:
             - build
```
