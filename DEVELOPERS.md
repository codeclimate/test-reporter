## Viewing Documentation Locally

```console
make all
for x in man/*.1; do man "$x"; done
```

## Sample Coverage Files

### JavaScript

- [Files](/examples/javascript)
- [Project](https://github.com/codeclimate/javascript-test-reporter/tree/cb80deb6667f62b701dcfea47b9143ceea6c7c1d)

### PHP

- [File](/examples/clover.xml)
- [Project](https://github.com/codeclimate/php-test-reporter/tree/e86f3e6105796dfc7a43b1eb7da5d1039388052c)

### Python

- [File](/examples/.coverage)
- [Project](https://github.com/codeclimate/python-test-reporter/tree/dc5236b37fe5eac5604bd7b1384d382072e4fd43)

### Ruby

- [Files](/examples/ruby)
- [Project](https://github.com/codeclimate/ruby-test-reporter/tree/1ec10f635414c70dc4e9c102825557b6510a8037)

### Cobertura

- [File](/examples/cobertura/coverage.xml)
- [Project](https://github.com/codeclimate-testing/cobertura-example/tree/d4ae1230498120ca6160343943298bcf7c6f202c)

## Coverage Payload

The coverage payload expected by Code Climate is defined canonically by
schema.json in the root of this repository. Examples can be found in the examples
folder.

## Usage Agent

All uploads should occur with a user agent of:

```
TestReporter/{VERSION} (Code Climate, Inc.)
```

## Server-Side Notes

### Writing

- We default `run_at` to `Time.now` when missing
- We ignore `committed_at` when it's `0` (after `to_i`), apparently
  - No idea what downstream impact this case has, I know we at least rely on
    this in the Extension to know if results are still able to be rendered
- We just skip reports if `ci_service.pull_request != "false"`, so good thing
  no-one writes that?
- `environment` doesn't seem required, and we store only the following keys:

  - `:test_framework`
  - `:pwd`
  - `:rails_root`
  - `:simplecov_root`
  - `:gem_version`

- From `ci_service` we store the following keys:

  - `:name`
  - `:build_url`
  - `:build_identifier`
  - `:pull_request`
  - `:branch`
  - `:commit_sha`

### Reading

#### Coverage Comparisons

Used for comparison and/or status update event payloads:

- `branch`
- `commit_sha`
- `committed_at`
- `covered_percent`

#### Rendering Coverage Info on codeclimate.com

Unclear. It think only `covered_percent` is copied to `snapshot` (repo-totals)
and `constant` (by source file) records.

#### Rendering Coverage Info on api.codeclimate.com

TODO

## Releasing

To release a new version,

* Update VERSION in Makefile
* Update CHANGELOG
* Run
  ```
    make release
  ```

This command will build the binaries for the given version and update the binaries
for _latest_.
