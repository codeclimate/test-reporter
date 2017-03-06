## Viewing Documentation Locally

```console
make all
man man/cc-test-reporter.1
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

## Coverage Payload

This is the payload currently expected by `codeclimate.com/test_reports`.

```json
{
  "ci_service": {
    "branch": "",
    "build_identifier": "",
    "build_url": "",
    "commit_sha": "",
    "committed_at": "",
    "name": "",
    "pull_request": "",
    "worker_id": ""
  },
  "covered_percent": 100,
  "covered_strength": 1,
  "environment": {
    "gem_version": "",
    "package_version": "",
    "pwd": "",
    "rails_root": "",
    "reporter_version": "",
    "simplecov_root": ""
  },
  "git": {
    "branch": "",
    "committed_at": 1234567,
    "head": "",
  },
  "line_counts": {
    "covered": 1,
    "missed": 1,
    "total": 1
  },
  "partial": false,
  "repo_token": "",
  "run_at": 1234567,
  "source_files": [
    {
      "blob_id": "",
      "coverage": [null, 0, 1],
      "covered_percent": 100,
      "covered_strength": 1,
      "line_counts": {
        "covered": 1,
        "missed": 1,
        "total": 1
      },
      "name": ""
    }
  ]
}
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
