% CC-TEST-REPORTER-ENV(1) User Manuals
% Code Climate <hello@codeclimate.com>
% February 2017

# PROLOG

This is a sub-command of **cc-test-reporter**(1).

# SYNOPSIS

**cc-test-reporter-env**

# DESCRIPTION

Infer and output information about the environment the reporter is running in.

# EXAMPLE OUTPUT

The output is formatted for use with **eval**(1):

    GIT_BRANCH=master
    GIT_COMMIT_SHA=594a20638eb9a758e2481c9ad2bdae121a1e03ed
    GIT_COMMITTED_AT=1488138087
    CI_NAME=circle-ci
    CI_BUILD_ID=7
    CI_BUILD_URL=https://circleci.com/gh/foo/bar/7

# INFERENCE RULES

Any values set explicitly in the environment are output as-is. Unset values are
inferred using the following rules. If no value can be inferred, an empty
variable will be present in the output. Clients are expected to check for this
and error accordingly if they require a value.

## GIT_BRANCH

If *./.git* exists, read **git rev-parse --abbrev-ref HEAD**. Otherwise, try the
following environment variables in order:

    APPVEYOR_REPO_BRANCH
    BRANCH_NAME
    BUILDKITE_BRANCH
    CIRCLE_BRANCH
    CI_BRANCH
    CI_BUILD_REF_NAME
    TRAVIS_BRANCH
    WERCKER_GIT_BRANCH

## GIT_COMMIT_SHA

If *./.git* exists, read **git log -1 --pretty=format:'%H'**. Otherwise, try the
following environment variables in order:

    APPVEYOR_REPO_COMMIT
    BUILDKITE_COMMIT
    CIRCLE_SHA1
    CI_BUILD_REF
    CI_BUILD_SHA
    CI_COMMIT
    CI_COMMIT_ID
    GIT_COMMIT
    WERCKER_GIT_COMMIT

## GIT_COMMITTED_AT

If *./.git* exists, read **git log -1 --pretty=format:'%ct'**. Otherwise, try
the following environment variables in order:

    CI_COMMITED_AT [sic]

## CI_NAME

Chosen based on the presence (and possibly value) of one the following
environment variables:

    APPVEYOR
    BUILDKITE
    CIRCLECI
    CI_NAME
    GITLAB_CI
    JENKINS_URL
    SEMAPHORE
    TDDIUM
    TRAVIS
    WERCKER

## CI_BUILD_ID

Chosen from the first of:

    APPVEYOR_BUILD_ID
    BUILDKITE_JOB_ID
    BUILD_NUMBER
    CIRCLE_BUILD_NUM
    CI_BUILD_ID
    CI_BUILD_NUMBER
    SEMAPHORE_BUILD_NUMBER
    TDDIUM_SESSION_ID
    TRAVIS_JOB_ID
    WERCKER_BUILD_ID

## CI_BUILD_URL

Chosen from the first of:

    APPVEYOR_API_URL
    BUILDKITE_BUILD_URL
    BUILD_URL
    CIRCLE_BUILD_NUM
    CI_BUILD_URL
    WERCKER_BUILD_URL
