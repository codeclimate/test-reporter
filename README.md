# cc-test-reporter

TODO: details

## Installation

```
curl -L https://codeclimate.s3.amazonaws.com/test-reporter/test-reporter-latest > ./cc-test-reporter
chmod +x ./cc-test-reporter
sudo mv ./cc-test-reporter /usr/local/bin # anywhere in $PATH
```

## Usage

See the [man-pages][man/].

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
   cc-test-reporter upload-coverage
   ```

## Copyright

See the [LICENSE][].
