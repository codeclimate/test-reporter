#!/bin/sh

set -o nounset
set -o errexit
set -o pipefail
            
# Install go
curl -O https://dl.google.com/go/go1.15.darwin-amd64.tar.gz
tar -xzf go1.15.darwin-amd64.tar.gz
echo 'export PATH=$PATH:$PWD/go/bin' >> "$BASH_ENV"

# Set go path
mkdir -p ~/gopath/src/github.com/codeclimate/test-reporter
echo 'export GOPATH=$HOME/gopath' >> "$BASH_ENV"
. "$BASH_ENV"
cd $GOPATH/src/github.com/codeclimate/test-reporter
cp -r ~/project/ $GOPATH/src/github.com/codeclimate/test-reporter/
