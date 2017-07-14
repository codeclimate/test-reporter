# Change log

## Unreleased

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
