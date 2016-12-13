#!/bin/bash -e

. bin/commons.sh

if [ -d "$GOPATH" ]; then
    tar -zcvf $GOPATH/bin/valkiria-${VERSION}.tar.gz $GOPATH/bin/valkiria
else
    echo "target file not available, please run 'make compile' first"
fi
