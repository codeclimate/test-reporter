# Change log

## v0.7.0 [(2019-07-31)](https://github.com/codeclimate/test-reporter/releases/tag/v0.7.0)

* [NEW] Add support to xccov JSON report (Swift/Xcode 11) [#399][]


## v0.6.4 [(2019-01-22)](https://github.com/codeclimate/test-reporter/releases/tag/v0.6.4)

* [NEW] Updates release strategy to build a new Linux binary that uses
  `netcgo` [#382][]

[#382]: https://github.com/codeclimate/test-reporter/pull/382

## v0.6.3 [(2018-09-04)](https://github.com/codeclimate/test-reporter/releases/tag/v0.6.3)

* [NEW] Updates release strategy to build a new Linux binary that uses
  `netcgo` [#355][]

[#355]: https://github.com/codeclimate/test-reporter/pull/355

## v0.6.2 [(2018-07-31)](https://github.com/codeclimate/test-reporter/releases/tag/v0.6.2)

* [NEW] Add support for multiple JaCoCo source paths [#348][]

## v0.6.1 [(2018-07-30)](https://github.com/codeclimate/test-reporter/releases/tag/v0.6.1)

* [FIX] Support clover `path` attribute on file nodes in XML report [#349][]

[#349]: https://github.com/codeclimate/test-reporter/pull/349

## v0.6.0 [(2018-05-21)](https://github.com/codeclimate/test-reporter/releases/tag/v0.6.0)

* [FIX] Update `Gcov` formatter to report the correct source file paths [#338][]

[#338]: https://github.com/codeclimate/test-reporter/pull/338

## v0.5.2 [(2018-05-15)](https://github.com/codeclimate/test-reporter/releases/tag/v0.5.2)

* [FIX] Update `Cobertura` formatter to ignore invalid line numbers in a
  `cobertura.xml` file [#335][]

[#335]: https://github.com/codeclimate/test-reporter/pull/335

## v0.5.1 [(2018-03-19)](https://github.com/codeclimate/test-reporter/releases/tag/v0.5.1)

* [FIX] Fix bug with formatting JaCoCo test coverage [#318][]

[#318]: https://github.com/codeclimate/test-reporter/pull/318

## v0.5.0 [(2018-02-26)](https://github.com/codeclimate/test-reporter/releases/tag/v0.5.0)

* [NEW] Add flag to upload coverage insecurely [#310][]

[#310]: https://github.com/codeclimate/test-reporter/pull/310

## v0.4.5 [(2018-02-19)](https://github.com/codeclimate/test-reporter/releases/tag/v0.4.5)

* [FIX] Add partial automated support for Heroku CI builds [#305][]

[#305]: https://github.com/codeclimate/test-reporter/pull/305

## v0.4.4 [(2018-02-15)](https://github.com/codeclimate/test-reporter/releases/tag/v0.4.4)

* [FIX] Add support for Codeship CI environment variables [#300][]

## v0.4.3 [(2018-01-18)](https://github.com/codeclimate/test-reporter/releases/tag/v0.4.3)

[#300]: https://github.com/codeclimate/test-reporter/pull/300

* [FIX] Fix logging to not include extraneous newline, correct spelling [#288][]

[#288]: https://github.com/codeclimate/test-reporter/pull/288

## v0.4.2 [(2018-01-09)](https://github.com/codeclimate/test-reporter/releases/tag/v0.4.2)

* [FIX] Improved performance of the Cobertura, Gcov, and Jacoco formatters [#285][]

[#285]: https://github.com/codeclimate/test-reporter/pull/285

## v0.4.1 [(2018-01-08)](https://github.com/codeclimate/test-reporter/releases/tag/v0.4.1)

* [FIX] Improved performance of the LCOV formatter [#270][]

[#270]: https://github.com/codeclimate/test-reporter/pull/270

## v0.4.0 [(2017-12-19)](https://github.com/codeclimate/test-reporter/releases/tag/v0.4.0)

* [NEW] Add support for `excoveralls` json report (Elixir) [#278][]

[#278]: https://github.com/codeclimate/test-reporter/pull/278

## v0.3.4 [(2017-12-13)](https://github.com/codeclimate/test-reporter/releases/tag/v0.3.4)

* [FIX] `sum-coverage` command when merging source files with
        different coverage length but same blob ID [#272][]

[#272]: https://github.com/codeclimate/test-reporter/pull/272

## v0.3.3 [(2017-12-05)](https://github.com/codeclimate/test-reporter/releases/tag/v0.3.3)

* [FIX] Treat 409 report upload status as warning and exit 0 [#268][]

[#268]: https://github.com/codeclimate/test-reporter/pull/268

## v0.3.2 [(2017-10-30)](https://github.com/codeclimate/test-reporter/releases/tag/v0.3.2)

* [FIX] Update `coverage.py` formatter to parse `<source>` tags and correctly create
  source file paths. [#247][]

[#247]: https://github.com/codeclimate/test-reporter/pull/247

## v0.3.1 [(2017-10-02)](https://github.com/codeclimate/test-reporter/releases/tag/v0.3.1)

* [FIX] "format-coverage" "--add-prefix" option when is an empty string

## v0.3.0 [(2017-10-02)](https://github.com/codeclimate/test-reporter/releases/tag/v0.3.0)

* [NEW] "format-coverage" now supports a new option "--add-prefix" for prefixing a path to all file paths
* [FIX] JaCoCo formatter now properly takes into account package when creating file paths
* [FIX] "format-coverage" "--prefix" option now accepts paths with or without trailing slash

## v0.2.2 [(2017-09-06)](https://github.com/codeclimate/test-reporter/releases/tag/v0.2.2)

* [FIX] Fix build task to statically compile binaries [#221][]

[#221]: https://github.com/codeclimate/test-reporter/pull/221

## v0.2.1 [(2017-09-05)](https://github.com/codeclimate/test-reporter/releases/tag/v0.2.1)

* [FIX] Update `Cobertura` formatter to parse `<source>` tags and correctly source file
   paths. [#218][]

[#218]: https://github.com/codeclimate/test-reporter/pull/218

## v0.2.0 [(2017-08-24)](https://github.com/codeclimate/test-reporter/releases/tag/v0.2.0)

* [NEW] Add support for DroneCI by @teohm [#215][]
* [FIX] Fix Clover formatter [#213][]
* [FIX] Raise an error when report is empty [#214][]

[#213]: https://github.com/codeclimate/test-reporter/pull/213
[#214]: https://github.com/codeclimate/test-reporter/pull/214
[#215]: https://github.com/codeclimate/test-reporter/pull/215

## v0.1.14 [(2017-08-02)](https://github.com/codeclimate/test-reporter/releases/tag/v0.1.14)

* [NEW] Add support for `SSL_CERT_FILE` env var on
  `after-build`/`upload-coverage` commands
* [FIX] Update `Cobertura` formatter to correctly parse coverage information
  when file contains inner classes.

## v0.1.13 [(2017-07-14)](https://github.com/codeclimate/test-reporter/releases/tag/v0.1.13)

* [FIX] Update `Cobertura` formatter to correctly parse coverage information
  from a `cobertura.xml` file that has lines not sorted by number or when lines
  with the same number are present more than one time.

## v0.1.12 [(2017-07-14)](https://github.com/codeclimate/test-reporter/releases/tag/v0.1.12)

* [FIX] `format-coverage`/`after-build` `--debug` outputs `codeclimate.json` content

## v0.1.11 [(2017-07-06)](https://github.com/codeclimate/test-reporter/releases/tag/v0.1.11)

* [NEW] Add `JaCoCo` support
* [FIX] `upload-coverage` outputs message when successful

## v0.1.10 [(2017-06-08)](https://github.com/codeclimate/test-reporter/releases/tag/v0.1.10)

* [NEW]  Add support on `format-coverage` for Travis ENV vars to
  infer correctly git commit sha and git branch name
* [FIX] `format-coverage` when `--input-type` is specified without a file path

## v0.1.9 [(2017-06-06)](https://github.com/codeclimate/test-reporter/releases/tag/v0.1.9)

* [NEW] Add `SwiftCov` support
* [NEW] Add `Cobertura` support
* [FIX] `sum-coverage` output to `stdout`

## v0.1.8 [(2017-05-21)](https://github.com/codeclimate/test-reporter/releases/tag/v0.1.8)

* [NEW] `sum-coverage` when merging source file coverage it preserves nulls

## v0.1.7 [(2017-05-20)](https://github.com/codeclimate/test-reporter/releases/tag/v0.1.7)

* [FIX] Raise an error when invalid format-coverage path usage
* [FIX] Avoid accessing the git repository for sum-coverage
* [FIX] Avoid accessing the git repository for upload-coverage
* [FIX] Improve performance of sum-coverage

## v0.1.6 [(2017-05-18)](https://github.com/codeclimate/test-reporter/releases/tag/v0.1.6)

* [FIX] `format-coverage` when COVERAGE_FILE arg is not present

## v0.1.5 [(2017-05-17)](https://github.com/codeclimate/test-reporter/releases/tag/v0.1.5)

* [NEW] Add coverage file path argument to `format-coverage`
* [NEW] Add debug logs to Git commands

## v0.1.4 [(2017-05-08)](https://github.com/codeclimate/test-reporter/releases/tag/v0.1.4)

* [NEW] Add `Gocov` support
* [NEW] Add `Clover XML` support
* [NEW] Add coverage strength for tools that don't calculate it.
* [FIX] Ensure "blank" lines on source file's coverage

## v0.1.3 [(2017-05-03)](https://github.com/codeclimate/test-reporter/releases/tag/v0.1.3)

* [NEW] Add `--parts`(`-p`) flag to `sum-coverage`
* [NEW] Add support to repositories without `.git` access
* [NEW] Add `after-build` command
* [NEW] Add support for `coverage.py`

## v0.1.2 [(2017-04-24)](https://github.com/codeclimate/test-reporter/releases/tag/v0.1.2)

* [NEW] Add `before-build` command
* [NEW] Add `--prefix`(`-p`) flag to `format-coverage`

## v0.1.1 [(2017-04-18)](https://github.com/codeclimate/test-reporter/releases/tag/v0.1.1)

* [NEW] Add `--input-type`(`-t`) flag to `format-coverage`
* [NEW] Add `--debug`(`-d`) flag
* [NEW] Add `lcov` support
* [FIX] Source file `blob_id`

## alpha [(2017-03-22)](https://github.com/codeclimate/test-reporter/releases/tag/alpha)

### New features

* [NEW] Add `format-coverage` command
* [NEW] Add `sum-coverage` command
* [NEW] Add `upload-coverage` command
* [NEW] Add `simplecov` support
