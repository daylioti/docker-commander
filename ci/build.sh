#!/usr/bin/env bash


GOARCH=${_GOARCH}
GOOS=${_GOOS}
/bin/bash ./ci/install_upx.sh
env GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags "-s -w -X github.com/daylioti/docker-commander/version.Version=${TRAVIS_BRANCH}" -o ${NAME} || ERROR=true
./upx --brute docker-commander
mkdir -p dist


FILE=${NAME}_${TRAVIS_BRANCH}_${GOOS}_${GOARCH}

tar -czf dist/${FILE}.tgz ${NAME} || ERROR=true

if [ "$ERROR" == "true" ]; then
        exit 1
fi

