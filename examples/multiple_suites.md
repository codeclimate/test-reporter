## Example 1
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
  podman:
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

