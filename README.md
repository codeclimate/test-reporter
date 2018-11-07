![code-climate-logo_black-200](https://user-images.githubusercontent.com/18341459/47682820-32937480-db93-11e8-9d81-e5052a22453b.png)

# Code Climate Test Reporter

Code Climate's test reporter is a binary that works in coordination with codeclimate.com to report test coverage data. Once you've set up test coverage reporting you can:
* view test coverage reports for each file alongside quality metrics like complexity, duplication, and churn,
* toggle between viewing code issues and test coverage line-by-line in the same source listings,
* block PRs from being merged if they don't meet your team's standards for test coverage percentage.

Code Climate accepts test coverage data from virtually any location, including locally run tests or your continuous integration (CI) service, and supports a variety of programming languages and test coverage formats, including Ruby, JavaScript, Go, Python, PHP, Java, and more.

For installation instructions, check out our docs on [Configuring Test Coverage](https://docs.codeclimate.com/docs/configuring-test-coverage) and [Test Coverage Troubleshooting Tips](https://docs.codeclimate.com/docs/test-coverage-troubleshooting-tips).

Some installations may require the use of the following [subcommands](https://docs.codeclimate.com/docs/configuring-test-coverage#section-list-of-subcommands): 



#### `format-coverage` - formats test report from local test suite into generalized format, readable by Code Climate

- `-t` or  `--input-type` *simplecov | lcov | coverage.py | gcov | clover* - Identifies the input type (format) of the COVERAGE_FILE

- `-o PATH` or  `--output PATH` - Output to PATH. If - is given, content will be written to stdout. Defaults to coverage/codeclimate.json.

- `-p PATH` or `--prefix PATH` - The prefix to remove from absolute paths in coverage payloads, to make them relative to the project root. This is usually the directory in which the tests were run. Defaults to current working directory.

- `COVERAGE_FILE` - Path to the coverage file to process. Defaults to searching known paths where coverage files could exist and selecting the first one found.



#### `sum-coverage` - combines test reports from multiple sources (i.e. multiple test suites or parallelized CI builds) into one test report which is readable by Code Climate

- `-o PATH` or  `--output PATH` - Output to PATH. If - is given, content will be written to stdout. Defaults to coverage/codeclimate.json.

- `-p NUMBER` or `--parts NUMBER` - Expect NUMBER payloads to sum. If this many arguments are not present, command will error. This ensures you don't accidentally sum incomplete results.







#### `upload-coverage` - uploads formatted, singular test report to Code Climate API

- `-i PATH` or `--input PATH` - Read payload from PATH. If - is given, the payload will be read from stdin. Defaults to coverage/codeclimate.json.

- `-r ID` or  `--id ID` - The reporter identifier to use when reporting coverage information. The appropriate value can be found in your Repository Settings page on codeclimate.com. Defaults to the value in the `CC_TEST_REPORTER_ID` environment variable. The uploader will error if a value is not found.

- `-e URL` or `--endpoint URL` - The endpoint to upload coverage information to. Defaults to the value in the CC_TEST_REPORTER_COVERAGE_ENDPOINT environment variable, or a hard-coded default (currently "https://codeclimate.com/test_reports").



#### `after-build` - combines `format-coverage` and `upload-coverage`

- `--exit-code $EXIT_CODE` - `$EXIT_CODE` should be the exit code of your test suite process. Some CI system expose this as an environment variable; for others, you may need to manually capture `$?` to provide it to `after-build` later. Providing this will prevent sending test coverage results for failed tests.


To sign up for Code Climate, head [here](https://codeclimate.com/quality/pricing/).


## Copyright

See the [LICENSE](https://github.com/codeclimate/test-reporter/blob/master/LICENSE).
