![code-climate-logo_black-200](https://user-images.githubusercontent.com/18341459/47682820-32937480-db93-11e8-9d81-e5052a22453b.png)

# Code Climate Test Reporter

Code Climate's test reporter is a binary that works in coordination with codeclimate.com to report test coverage data. Once you've set up test coverage reporting you can:
* view test coverage reports for each file alongside quality metrics like complexity, duplication, and churn,
* toggle between viewing code issues and test coverage line-by-line in the same source listings,
* block PRs from being merged if they don't meet your team's standards for test coverage percentage.

Code Climate accepts test coverage data from virtually any location, including locally run tests or your continuous integration (CI) service, and supports a variety of programming languages and test coverage formats, including Ruby, JavaScript, Go, Python, PHP, Java, and more.

For installation instructions, check out our docs on [Configuring Test Coverage](https://docs.codeclimate.com/docs/configuring-test-coverage) and [Test Coverage Troubleshooting Tips](https://docs.codeclimate.com/docs/test-coverage-troubleshooting-tips).

To sign up for Code Climate, head [here](https://codeclimate.com/quality/pricing/).

# Versioning
The test reporter's current version is documented in [VERSIONING/VERSION](https://github.com/codeclimate/test-reporter/blob/master/VERSIONING/VERSION), following the [Semantic Versioning](https://semver.org/) convention.

# Binaries

## Download
The test reporter is distributed as a pre-built binary named cc-test-reporter. You can fetch the pre-built binary from the following URLs:

### Linux
- [codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64](https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64)
- [codeclimate.com/downloads/test-reporter/test-reporter-${VERSION}-linux-amd64](https://codeclimate.com/downloads/test-reporter/test-reporter-${VERSION}-linux-amd64)


### Linux netcgo (recommended if you're using a VPN)
- [codeclimate.com/downloads/test-reporter/test-reporter-latest-netcgo-linux-amd64](https://codeclimate.com/downloads/test-reporter/test-reporter-latest-netcgo-linux-amd64)
- [codeclimate.com/downloads/test-reporter/test-reporter-${VERSION}-netcgo-linux-amd64](https://codeclimate.com/downloads/test-reporter/test-reporter-${VERSION}-netcgo-linux-amd64)

### OS X
- [codeclimate.com/downloads/test-reporter/test-reporter-latest-darwin-amd64](https://codeclimate.com/downloads/test-reporter/test-reporter-latest-darwin-amd64)
- [codeclimate.com/downloads/test-reporter/test-reporter-${VERSION}-darwin-amd64](https://codeclimate.com/downloads/test-reporter/test-reporter-${VERSION}-darwin-amd64)


#### e.g
>```console
>$ curl -L -O https://codeclimate.com/downloads/test-reporter/test-reporter-0.10.1-darwin-amd64
>```

## Verifying binaries

Along with the binaries you can download a file with a SHA 256 checksum for the given version from the link shown below, or you can attach it to your clipboard from the [docs page](https://docs.codeclimate.com/docs/configuring-test-coverage#locations-of-pre-built-binaries).

- [codeclimate.com/downloads/test-reporter/test-reporter-${VERSION}-darwin-amd64.sha256](https://codeclimate.com/downloads/test-reporter/test-reporter-${VERSION}-darwin-amd64.sha256)

To download the file containing the checksum using `curl`:
#### e.g
>```console
>$ curl -L -O https://codeclimate.com/downloads/test-reporter/test-reporter-0.10.1-darwin-amd64.sha256
>```

To check that a downloaded file matches the checksum, run it through `shasum` with a command such as:

```console
$ grep test-reporter-${VERSION}-darwin-amd64  test-reporter-${VERSION}-darwin-amd64.sha256 | shasum -a 256 -c -
```

The GPG detached signature of SHA checksums can be download analogously from the following url:

- [codeclimate.com/downloads/test-reporter/test-reporter-${VERSION}-darwin-amd64.sha256.sig](https://codeclimate.com/downloads/test-reporter/test-reporter-${VERSION}-darwin-amd64.sha256.sig)

You can use it with `gpg` to verify the integrity of your downloaded checksum. You will first need to import
the GPG publick key. To import the key:

```console
$ gpg --keyserver  keys.openpgp.org --recv-keys 9BD9E2DD46DA965A537E5B0A5CBF320243B6FD85
```

Then use the following command to verify the file's signature.

```console
$ gpg --verify test-reporter-${VERSION}-darwin-amd64.sha256.sig test-reporter-${VERSION}-darwin-amd64.sha256
```

## Copyright

See the [LICENSE](https://github.com/codeclimate/test-reporter/blob/master/LICENSE).
