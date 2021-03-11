![code-climate-logo_black-200](https://user-images.githubusercontent.com/18341459/47682820-32937480-db93-11e8-9d81-e5052a22453b.png)

# Code Climate Test Reporter

Code Climate's test reporter is a binary that works in coordination with codeclimate.com to report test coverage data. Once you've set up test coverage reporting you can:
* view test coverage reports for each file alongside quality metrics like complexity, duplication, and churn,
* toggle between viewing code issues and test coverage line-by-line in the same source listings,
* block PRs from being merged if they don't meet your team's standards for test coverage percentage.

Code Climate accepts test coverage data from virtually any location, including locally run tests or your continuous integration (CI) service, and supports a variety of programming languages and test coverage formats, including Ruby, JavaScript, Go, Python, PHP, Java, and more.

For installation instructions, check out our docs on [Configuring Test Coverage](https://docs.codeclimate.com/docs/configuring-test-coverage) and [Test Coverage Troubleshooting Tips](https://docs.codeclimate.com/docs/test-coverage-troubleshooting-tips).

To sign up for Code Climate, head [here](https://codeclimate.com/quality/pricing/).

## Releasing a new version

Test reporter's new versions are released automatically when pushing to branches that match `vx.x.x`. The test reporter's current version is documented in [VERSIONING/VERSION](https://github.com/codeclimate/test-reporter/blob/master/VERSIONING/VERSION), following the [Semantic Versioning](https://semver.org/) convention.

There are two script helpers for creating a new release:
- [release-scripts/prep-release](https://github.com/codeclimate/test-reporter/blob/master/release-scripts/prep-release) which will create a new pull request, patching the current version. If you need to create a new MINOR or MAJOR version creating a manual pull request is the way to go.
- [release-scripts/release](https://github.com/codeclimate/test-reporter/blob/master/release-scripts/release) This script will create a new branch named `vx.x.x` that matches the version indicated in [VERSIONING/VERSION](https://github.com/codeclimate/test-reporter/blob/master/VERSIONING/VERSION), which should trigger the CI for creating a new release.

## Copyright

See the [LICENSE](https://github.com/codeclimate/test-reporter/blob/master/LICENSE).
