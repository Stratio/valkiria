#!/bin/bash -xv

VIRTUALENV=$PWD/target
ALLPACKAGE="./valkiria ./routes ./proc ./dbus"
LOGLEVEL=INFO
VERSION=$(cat VERSION)
BUILD=$(git rev-parse HEAD)
LDFLAGS="-X main.Version=$VERSION -X main.Build=$BUILD -X main.LOGLEVEL=$LOGLEVEL"

mkdir -p $VIRTUALENV/bin $VIRTUALENV/pkg $VIRTUALENV/src/github.com/Stratio/valkiria
export GOPATH=$VIRTUALENV
go get -u github.com/tools/godep
ln -s $PWD/valkiria $VIRTUALENV/src/github.com/Stratio/valkiria/valkiria
ln -s $PWD/Godeps $VIRTUALENV/src/github.com/Stratio/valkiria/Godeps
ln -s $PWD/dbus $VIRTUALENV/src/github.com/Stratio/valkiria/dbus
ln -s $PWD/routes /$VIRTUALENV/src/github.com/Stratio/valkiria/routes
ln -s $PWD/proc $VIRTUALENV/src/github.com/Stratio/valkiria/proc
ln -s $PWD/vendor $VIRTUALENV/src/github.com/Stratio/valkiria/vendor
cd $GOPATH/src/github.com/Stratio/valkiria
$GOPATH/bin/godep go install -v -ldflags "-w $LDFLAGS" $ALLPACKAGE