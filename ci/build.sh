#!/usr/bin/env bash


GOARCH=${_GOARCH}
GOOS=${_GOOS}

env GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${NAME} || ERROR=true

mkdir -p dist


FILE=${NAME}_${TRAVIS_BRANCH}_${GOOS}_${GOARCH}

tar -czf dist/${FILE}.tgz ${NAME} || ERROR=true

if [ "$ERROR" == "true" ]; then
        exit 1
fi

