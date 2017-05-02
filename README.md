# cc-test-reporter

Report information about your CI builds to Code Climate.

## Installation & Usage

Most CI systems allow configuration of commands to run as part of setup, before,
and after a test build. Using Circle CI as an example:

```yaml
machine:
  environment:
    CC_TEST_REPORTER_ID: ...

dependencies:
  post:
    - curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter
    - chmod +x ./cc-test-reporter

test:
  pre:
    - ./cc-test-reporter before-build

  post:
    - ./cc-test-reporter after-build --exit-code $EXIT_CODE
```

  Where:

  - `$EXIT_CODE` should be the exit code of your test suite process. Some CI
     system expose this as an environment variable; for others, you may need
     to manually capture $? to provide it to after-build later. Providing this
     will prevent sending test coverage results for failed tests. 

## Code Climate: Enterprise

To report coverage to your locally-hosted *Code Climate: Enterprise* instance,
export the `CC_TEST_REPORTER_COVERAGE_ENDPOINT` variable, or pass the
`--coverage-endpoint` option to `after-build`.

```sh
CC_TEST_REPORTER_COVERAGE_ENDPOINT=https://codeclimate.my-domain.com/test_reports
```

## Low-level Usage

The test reporter is implemented as a composition of lower-level commands, which
may themselves be useful. See the [man-pages](man) for details of these
commands.

## Parallel Tests

Code Climate supports parallel test setups using sub-commands provided by the
test reporter.  Specifically, the test reporter has sub-commands to:

1. format partial results (`format-coverage`)
1. sum partial results into a single result (`sum-coverage`) and
1. upload the single result to Code Climate (`upload-coverage`)

To make use of these commands, parallel test support requires:

1. the ability to run commands after *each* batch of tests has completed (most CI systems support this)
1. the ability to run commands after *all* tests have completed (most CI systems support this)
1. uploading and downloading partial test coverage data to/from shared storage
(using AWS S3, for example)

For example:

1. After *each* batch of tests:

   ```sh
   ./cc-test-reporter format-coverage --output "coverage/codeclimate.$N.json"
   aws s3 sync coverage/ "s3://my-bucket/coverage/$SHA"
   ```

   Where:

   - `$N` should be a unique identifier for that batch of tests
   - `$SHA` should be the commit for which the coverage was generated; you can
     use an existing, CI-provided variable or `./cc-test-reporter env` to infer
     `$GIT_COMMIT_SHA` and use that.

1. After *all* tests:

   ```sh
   aws s3 sync "s3://my-bucket/coverage/$SHA" coverage/
   cc-test-reporter sum-coverage --output - --parts $PARTS coverage/codeclimate.*.json | \
     cc-test-reporter upload-coverage --input -
   ```

   Where:

   - `$PARTS` should be the number of payloads to sum.

## Multiple Suites

Coverage from multiple suites can be sent to Code Climate by aggregating each
suite's results into one final report.

1. After each test suite, run:

  ```sh
  ./cc-test-reporter format-coverage --output coverage/codeclimate.$SUITE.json
  ```

1. After all test suites, run:

  ```sh
  ./cc-test-reporter sum-coverage coverage/codeclimate.*.json | \
    ./cc-test-reporter upload-coverage
  ```

## Copyright

See the [LICENSE](LICENSE).
