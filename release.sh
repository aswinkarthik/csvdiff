#!/bin/bash

set -e

sudo apt-get update -y && sudo apt-get install rpm -y
test -n "$TRAVIS_TAG" && curl -sL https://git.io/goreleaser | bash