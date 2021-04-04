#!/usr/bin/env bash


GOARCH=${_GOARCH}
GOOS=${_GOOS}
env GOOS=${GOOS} GOARCH=${GOARCH} go build -mod=mod -ldflags "-s -w -X github.com/daylioti/docker-commander/version.Version=${TRAVIS_BRANCH}" -o ${NAME} || ERROR=true
if [ "${_GOOS}" == "linux" ]; then
  /bin/bash ./ci/install_upx.sh
  ./upx --brute docker-commander
fi
mkdir -p dist


FILE=${NAME}_${TRAVIS_BRANCH}_${GOOS}_${GOARCH}

tar -czf dist/${FILE}.tgz ${NAME} || ERROR=true

if [ "$ERROR" == "true" ]; then
        exit 1
fi

