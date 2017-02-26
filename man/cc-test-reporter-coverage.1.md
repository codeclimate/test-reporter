% CC-TEST-REPORTER-COVERAGE(1) User Manuals
% Code Climate <hello@codeclimate.com>
% February 2017

# PROLOG

This is a sub-command of **cc-test-reporter**(1).

# SYNOPSIS

**cc-test-reporter-coverage**

# DESCRIPTION

Format and upload coverage information from supported sources to Code Climate
servers.

This is roughly equivalent to:

    eval $(cc-test-reporter env)
    cc-test-reporter format-coverage
    cc-test-reporter upload-coverage

# ENVIRONMENT VARIABLES

*CC_TEST_REPORTER_TOKEN*

# SEE ALSO

**cc-test-reporter-env**(1),
**cc-test-reporter-format-coverage**(1), and
**cc-test-reporter-upload-coverage**(1).
