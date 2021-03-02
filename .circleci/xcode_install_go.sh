#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

PROJECT_RELATIVE_PATH=src/github.com/codeclimate/test-reporter
            
# Install go
curl -O https://dl.google.com/go/go1.15.darwin-amd64.tar.gz
tar -xzf go1.15.darwin-amd64.tar.gz
echo 'export PATH=$PATH:$PWD/go/bin' >> "$BASH_ENV"

# Set go path
mkdir -p ~/gopath/${PROJECT_RELATIVE_PATH}
echo 'export GOPATH=$HOME/gopath' >> "$BASH_ENV"
. "$BASH_ENV"
cd $GOPATH/${PROJECT_RELATIVE_PATH}
cp -r ~/project/ $GOPATH/${PROJECT_RELATIVE_PATH}
