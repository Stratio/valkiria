#!/bin/bash -e

VIRTUALENV=$PWD/target
export GOPATH=$VIRTUALENV
tar -zcvf valkiria.tar.gz $GOPATH/bin/valkiria