% CC-TEST-REPORTER-UPLOAD-COVERAGE(1) User Manuals
% Code Climate <hello@codeclimate.com>
% February 2017

# PROLOG

This is a sub-command of **cc-test-reporter**(1).

# SYNOPSIS

**cc-test-reporter-upload-coverage** [--input=\<path>]

# DESCRIPTION

Aggregate and upload formatted coverage payloads to Code Climate servers.

# OPTIONS

## -i, --input *PATH*

Read payload(s) from *PATH*. If a directory is given, payloads will be read from
*PATH*/\*.json. If *-* is given, a single payload will be expected on *stdin*.
Defaults to *coverage/*, a directory.

# ENVIRONMENT VARIABLES

*CC_TEST_REPORTER_TOKEN* or *CODECLIMATE_REPO_TOKEN* (deprecated).

# SEE ALSO

**cc-test-reporter-format-coverage**(1).
