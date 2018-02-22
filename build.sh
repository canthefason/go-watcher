#!/bin/bash -xe
#if [ "$TRAVIS_BRANCH" == "master" ]; then
  docker login -u "$DOCKER_USERNAME" -p "$DOCKER_PASSWORD"
  GO_VERSION=$(go version | awk '{print $3}')

  docker build -f Dockerfiles/Dockerfile.$GO_VERSION -t canthefason/go-watcher:$WATCHER_VERSION-$GO_VERSION .

  docker push canthefason/go-watcher:$WATCHER_VERSION-$GO_VERSION
#fi
