#!/bin/sh

set -o nounset
set -o errexit
set -o pipefail

ARCHS="linux-amd64 netcgo-linux-amd64 darwin-amd64"
S3_BASE_URL="https://s3.amazonaws.com/codeclimate/test-reporter"
ARTIFACTS_OUTPUT=artifacts.tar.gz

# Download artifacts from AWS
for arch in ${ARCHS}
do
    curl ${S3_BASE_URL}/test-reporter-${1}-${arch} -O
done

tar -c -f ${ARTIFACTS_OUTPUT} test-reporter-${1}-*

GITHUB_TOKEN=${GITHUB_TOKEN} hub release create -a ${ARTIFACTS_OUTPUT} -m "v${1}" ${1}
