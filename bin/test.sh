#!/bin/bash

VIRTUALENV=$PWD/target
PACKAGE="./valkiria ./routes ./proc ./dbus"

export GOPATH=$VIRTUALENV
cd $GOPATH/src/github.com/Stratio/valkiria
$GOPATH/bin/godep go test ${PACKAGE}
