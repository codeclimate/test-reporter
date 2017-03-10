% CC-TEST-REPORTER-UPLOAD-COVERAGE(1) User Manuals
% Code Climate <hello@codeclimate.com>
% February 2017

# PROLOG

This is a sub-command of **cc-test-reporter**(1).

# SYNOPSIS

**cc-test-reporter-upload-coverage** [--input=\<path>] [--id=\<id>] [--endpoint=\<url\>]

# DESCRIPTION

Upload pre-formatted coverage payloads to Code Climate servers.

# OPTIONS

## -i, --input *PATH*

Read payload from *PATH*. If *-* is given, the payload will be read from
*stdin*. Defaults to *coverage/codeclimate.json*.

## -r, --id *ID*

The reporter identifier to use when reporting coverage information. The
appropriate value can be found in your Repository Settings page on
*codeclimate.com*. Defaults to the value in the **CC_TEST_REPORTER_ID**
environment variable.

The uploader will error if a value is not found.

## -e, --endpoint *URL*

The endpoint to upload coverage information to. Defaults to the value in the
*CC_TEST_REPORTER_COVERAGE_ENDPOINT* environment variable, or a hard-coded
default (currently *"https://codeclimate.com/test_reports"*).

# ENVIRONMENT VARIABLES

*CC_TEST_REPORTER_ID*, *CC_TEST_REPORTER_COVERAGE_ENDPOINT*

The API endpoint to use
