#!/bin/bash -xv

VIRTUALENV=/tmp/path
ALLPACKAGE=./...
LOGLEVEL=INFO
VERSION=$(cat Version)
BUILD=$(git rev-parse HEAD)
LDFLAGS="-X main.Version=$VERSION -X main.Build=$BUILD -X main.LOGLEVEL=$LOGLEVEL"

mkdir -p $VIRTUALENV
export GOPATH=$VIRTUALENV
go get -u github.com/tools/godep
go get -u github.com/Stratio/valkiria
cd $GOPATH/src/github.com/Stratio/valkiria
$GOPATH/bin/godep go install -v -ldflags "-w $LDFLAGS" $ALLPACKAGE