#!/bin/bash

. bin/commons.sh

if [ -d "$GOPATH" ]; then
    cd $GOPATH/src/github.com/Stratio/valkiria
    # Publish coverage info in coveralls
    $GOPATH/bin/goveralls -coverprofile coverageAll.out -service jenkins
else
    echo "target file not available, please run 'make compile' first"
fi
