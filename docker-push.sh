#!/bin/bash -el

set -e

VERSION=${TRAVIS_TAG}
REPO_NAME="csvdiff"
GROUP="aswinkarthik"

docker build -t ${REPO_NAME}:${VERSION} .

docker tag ${REPO_NAME}:${VERSION} ${GROUP}/${REPO_NAME}:latest
docker tag ${REPO_NAME}:${VERSION} ${GROUP}/${REPO_NAME}:${VERSION}

docker push ${GROUP}/${REPO_NAME}:latest
docker push ${GROUP}/${REPO_NAME}:${VERSION}
