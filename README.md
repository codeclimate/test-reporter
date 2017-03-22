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

Code Climate doesn't yet support receiving partial payloads from (e.g.) parallel
test runs and summing them together server-side. However, the reporter does
support summing them together client-side.

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

## Copyright

See the [LICENSE](LICENSE).
