#!/bin/bash -e

VIRTUALENV=$PWD/target
export GOPATH=$VIRTUALENV
VERSION=$(cat VERSION)
tar -zcvf $GOPATH/bin/valkiria-${VERSION}.tar.gz $GOPATH/bin/valkiria
