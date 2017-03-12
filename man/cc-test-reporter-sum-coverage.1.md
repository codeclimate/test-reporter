% CC-TEST-REPORTER-SUM-COVERAGE(1) User Manuals
% Code Climate <hello@codeclimate.com>
% February 2017

# PROLOG

This is a sub-command of **cc-test-reporter**(1).

# SYNOPSIS

**cc-test-reporter-sum-coverage** [--output=\<path>] FILE [FILE, ...]

# DESCRIPTION

Combine (sum) multiple pre-formatted coverage payloads into one.

# OPTIONS

## -o, --output *PATH*

Output to *PATH*. If *-* is given, content will be written to *stdout*. Defaults
to *coverage/codeclimate.json*.

## FILE [FILE, ...]

Input files to combine. These are expected to be pre-formatted coverage
payloads. Passing a single file will return it unprocessed.

# INPUT VALIDATION

The following must be true for payloads to be sum-able. If these conditions are
not met, an error will be returned:

1. The value for *git.head* is equal across all payloads
1. All *source_files[].coverage* arrays for the same *name* are the same length

# ENVIRONMENT VARIABLES

None

# SEE ALSO

**cc-test-reporter-format-coverage**(1).
