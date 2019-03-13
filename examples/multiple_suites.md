## Example 1
- Multi-language
- Languages:
  - Ruby
  - Javascript
- CI: Circle 2.0
- Coverage Tools: 
  - Simplecov
  - lcov
- File: .circleci/config.yml
- OSS Repo: https://github.com/ale7714/loot

```
version: 2.0
defaults: &defaults
  working_directory: ~/repo
  docker:
    - image: circleci/ruby:2.4.2-jessie-node-browsers
      environment:
        PGHOST: 127.0.0.1
        PGUSER: loot_user
        RAILS_ENV: test
    - image: circleci/postgres:9.5-alpine
      environment:
        POSTGRES_USER: loot_user
        POSTGRES_DB: loot_test
        POSTGRES_PASSWORD: ""
jobs:
  build:
    <<: *defaults
    steps:
      - run:
          name:  Download cc-test-reporter
          command: |
            mkdir -p tmp/
            curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./tmp/cc-test-reporter
            chmod +x ./tmp/cc-test-reporter
      - persist_to_workspace:
          root: tmp
          paths:
            - cc-test-reporter
  backend-tests:
    <<: *defaults
    steps:
      - checkout
      - attach_workspace:
          at: ~/repo/tmp
      - run:
          name: Setup dependencies
          command: |
            bundle install --without production --path=${BUNDLE_PATH:-vendor/bundle}
            bundle exec rake db:create db:migrate
      - run:
          name: Run backend tests
          command: |
            bundle exec rake
            ./tmp/cc-test-reporter format-coverage -t simplecov -o tmp/codeclimate.backend.json coverage/backend/.resultset.json
      - persist_to_workspace:
          root: tmp
          paths:
            - codeclimate.backend.json
  frontend-tests:
    <<: *defaults
    steps:
      - checkout
      - attach_workspace:
          at: ~/repo/tmp
      - run: npm install
      - run:
          name: Run frontend testss
          command: |
            npm test
            ./tmp/cc-test-reporter format-coverage -t lcov -o tmp/codeclimate.frontend.json coverage/frontend/lcov.info
      - persist_to_workspace:
          root: tmp
          paths:
            - codeclimate.frontend.json
  upload-coverage:
    <<: *defaults
    environment:
      - CC_TEST_REPORTER_ID: 1acf55093f33b525eefdd9fb1e601d748e5d8b1267729176605edb4b5d82dc3d
    steps:
      - attach_workspace:
          at: ~/repo/tmp
      - run:
          name: Upload coverage results to Code Climate
          command: |
            ./tmp/cc-test-reporter sum-coverage tmp/codeclimate.*.json -p 2 -o tmp/codeclimate.total.json
            ./tmp/cc-test-reporter upload-coverage -i tmp/codeclimate.total.json
workflows:
  version: 2

  commit:
    jobs:
      - build
      - backend-tests:
          requires:
            - build
      - frontend-tests:
          requires:
             - build
      - upload-coverage:
          requires:
             - backend-tests
             - frontend-tests
```

## Example 2
- Parallel processes
- Languages:
  - Ruby
- CI: Travis CI
- Coverage Tools: 
  - Simplecov
- File: .travis.yml

```
language: ruby
rvm: "2.5.1"
sudo: true
filter_secrets: false

before_install:
  - 'curl "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" -o "awscli-bundle.zip"'
  - 'unzip awscli-bundle.zip'
  - './awscli-bundle/install -b ~/bin/aws'
  - 'export PATH=~/bin:$PATH'
  - gem update --system
  - gem install bundler -v 1.16.1
  - curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter
  - chmod +x ./cc-test-reporter

before_script:
  - bin/setup
  - ./cc-test-reporter before-build

script:
  - "bundle exec rake knapsack:rspec"
  - "bundle exec rake assets:precompile"

branches:
  only:
    - develop
    - master

cache:
  bundler: true
  yarn: true
  directories:
    - node_modules

env:
  global:
    - CC_TEST_REPORTER_ID=id
    - CI_NODE_TOTAL=5
    - AWS_ACCESS_KEY_ID=foo
    - AWS_SECRET_ACCESS_KEY=bar
    - AWS_DEFAULT_REGION=us-east-1
  matrix:
    - CI_NODE_INDEX=0
    - CI_NODE_INDEX=1
    - CI_NODE_INDEX=2
    - CI_NODE_INDEX=3
    - CI_NODE_INDEX=4

addons:
    postgresql: "9.6"
    elasticsearch: "5.x"
    chrome: stable

services:
  - postgresql
  - redis-server
  - elasticsearch

# Pipe the coverage data to Code Climate
after_script:
  - if [[ "$TRAVIS_TEST_RESULT" == 0 ]]; then ./cc-test-reporter format-coverage -t simplecov -o ./coverage/codeclimate.$CI_NODE_INDEX.json ./coverage/spec/.resultset.json; fi
  - if [[ "$TRAVIS_TEST_RESULT" == 0 ]]; then aws s3 sync coverage/ "s3://s3-bucket/coverage/$TRAVIS_BUILD_NUMBER"; fi
  - if [[ "$TRAVIS_TEST_RESULT" == 0 ]]; then aws s3 sync "s3://s3-bucket/coverage/$TRAVIS_BUILD_NUMBER" coverage/; fi
  - if [[ "$TRAVIS_TEST_RESULT" == 0 ]]; then ./cc-test-reporter sum-coverage --output - --parts $CI_NODE_TOTAL coverage/codeclimate.*.json | ./cc-test-reporter upload-coverage --input -; fi
```
