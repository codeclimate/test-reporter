# cc-test-reporter

Report information about your CI builds to Code Climate.

## Installation & Usage

**NOTE**: this is *proposed* usage, we've not built this yet.

Most CI systems allow configuration of commands to run as part of setup, before,
and after a test build. Using Circle CI as an example:

```yaml
machine:
  environment:
    CC_TEST_REPORTER_ID: ...

dependencies:
  post:
    - curl -L https://tools.codeclimate.com/test-reporter/test-reporter-latest > ./cc-test-reporter
    - chmod +x ./cc-test-reporter

test:
  pre:
    - ./cc-test-reporter before-build

  post:
    - ./cc-test-reporter after-build --exit-code $CIRCLE_EXIT_CODE
```

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

The test reporter supports aggregating multiple test coverage reports from, for example, parallel runs
together into a single report. After aggregation, the test reporter can send up the single aggregated
report to Code Climate.

This requires you store the partial payloads yourself after each test, then
download them to one location before using the reporter to upload a summed
payload.

For example:

1. After *each* test:

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
   cc-test-reporter sum-coverage --output - coverage/codeclimate.*.json | \
     cc-test-reporter upload-coverage --input -
   ```

## Multiple Suites

Coverage from multiple suites can be sent to Code Climate by aggregating each suite's results into one final report. To accomplish this, follow the instructions for [Parallel Tests](#parallel-tests) above. If the suites produce reports in the same location on disk, you can omit any steps that involve uploading and downloading partial reports.

## Copyright

See the [LICENSE](LICENSE).
