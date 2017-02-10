#!/bin/bash -e

. bin/commons.sh

if [ -d "$GOPATH" ]; then
    tar -cvf $GOPATH/bin/valkiria-${VERSION}.tar.gz -C $GOPATH/bin valkiria
else
    echo "target file not available, please run 'make compile' first"
fi
