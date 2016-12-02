#!/bin/bash -e

VIRTUALENV=$PWD/target

if [ -d "$VIRTUALENV" ]; then
    export GOPATH=$VIRTUALENV
    VERSION=$(cat VERSION)
    tar -zcvf $GOPATH/bin/valkiria-${VERSION}.tar.gz $GOPATH/bin/valkiria
else
    echo "target file not available, please run 'make compile' first"
fi
