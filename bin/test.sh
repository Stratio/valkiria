#!/bin/bash

VIRTUALENV=$PWD/target
PACKAGE="./valkiria ./routes ./proc ./dbus"

if [ -d "$VIRTUALENV" ]; then
    export GOPATH=$VIRTUALENV
    cd $GOPATH/src/github.com/Stratio/valkiria
    $GOPATH/bin/godep go test -v ${PACKAGE}
else
    echo "target file not available, please run 'make compile' first"
fi
