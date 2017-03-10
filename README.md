# cc-test-reporter

TODO: details

## Common Usage

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

*TODO*: the above is DRAFT content.

## Low-level Usage

See the [man-pages](man).

## Client-Side Aggregation (i.e Parallel Test Coverage)

TODO: describe further

1. After each test:

   ```
   eval $(cc-test-reporter env)
   cc-test-reporter format-coverage --out "coverage/codeclimate.$N.json"
   aws s3 sync coverage/ "s3://my-bucket/coverage/$GIT_COMMIT_SHA"
   ```

1. After all tests:

   ```
   eval $(cc-test-reporter env)
   aws s3 sync "s3://my-bucket/coverage/$GIT_COMMIT_SHA" coverage/
   cc-test-reporter sum-coverage --output - coverage/codeclimate.*.json | \
     cc-test-reporter upload-coverage --input -
   ```

## Copyright

See the [LICENSE](LICENSE).
