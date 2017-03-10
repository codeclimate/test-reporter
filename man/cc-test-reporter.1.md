% CC-TEST-REPORTER(1) User Manuals
% Code Climate <hello@codeclimate.com>
% February 2017

# SYNOPSIS

**cc-test-reporter** *COMMAND* [*OPTIONS*]

# DESCRIPTION

Determine and report information about your tests to Code Climate, so it can be
correlated with other analysis results.

# OPTIONS

## -v, --version

Print version information.

## -h, --help

Print the synopsis and list of commands. If used with a command, print help
information about that command.

## -d, --debug

Output debug messages during operation.

# COMMANDS

The reporter exposes high and low-level commands. For more details, see their
individual man-pages.

# HIGH-LEVEL COMMANDS

TODO:

# LOW-LEVEL COMMANDS

## cc-test-reporter-env(1)

Infer and print information about the environment where the reporter is running.

## cc-test-reporter-format-coverage(1)

Format coverage information from supported sources into a valid payload.

## cc-test-reporter-sum-coverage(1)

Combine multiple coverage payloads into one.

## cc-test-reporter-upload-coverage(1)

Upload formatted payloads to Code Climate servers.
