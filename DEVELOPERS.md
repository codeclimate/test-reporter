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
  ci_service: {
    branch:,
    build_identifier:,
    build_url:,
    commit_sha:,
    committed_at:,
    name:,
    pull_request:,
    worker_id:
  },
  covered_percent:,
  covered_strength:,
  environment: {
    gem_version:,
    pwd:,
    rails_root:,
    simplecov_root:
  },
  git: {
    branch:,
    committed_at:,
    head:
  },
  line_counts: {
    covered:,
    missed:,
    total:
  },
  //partial:,
  //repo_token:,
  run_at:,
  source_files: [
    {
      blob_id:,
      coverage: [
        // hit count, or null for missed
        ...,
        ...,
        ...,
      ],
      covered_percent:,
      covered_strength:,
      line_counts: {
        covered:,
        missed:,
        total:
      },
      name: name
    }
  ]
}
```
