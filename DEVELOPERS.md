## Viewing Documentation Locally

```console
make all
man man/cc-test-reporter.1
```

## Coverage Payload

This is the payload currently expected by `codeclimate.com/test_reports`.

*TODO*: remove keys not actually used by our system.

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
    "pwd": "",
    "rails_root": "",
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
      "coverage": [nil, 0, 1],
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
