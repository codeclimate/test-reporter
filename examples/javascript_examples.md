## Example 1
- Language: JavaScript
- CI: TravisCI (Environment Variables are set on travis web)
- Coverage Tool: [nyc](https://github.com/istanbuljs/nyc)
- File: travis.yml
- Single/Parallel: Single
- OSS Repo: https://github.com/amadeu01/adminviu

```yml
language: node_js

node_js:
  - 8.2.1
dist: trusty
addons:
  chrome: stable
before_script:
  - yarn global add nyc
  - curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter
  - chmod +x ./cc-test-reporter
  - ./cc-test-reporter before-build --debug

script:
  - nyc --reporter=lcov yarn run unit

after_script:
  - ./cc-test-reporter after-build --exit-code $TRAVIS_TEST_RESULT

notifications:
  email: false
```
