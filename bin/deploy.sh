#!/bin/bash -e

VIRTUALENV=$PWD/target

if [ -d "$VIRTUALENV" ]; then
    BASEDIR=`dirname $0`/..
    VERSION=`cat $BASEDIR/VERSION`
    # Upload normal universe
    cd ${VIRTUALENV}/bin && curl -sS -u stratio:${NEXUSPASS} --upload-file valkiria-${VERSION}.tar.gz http://sodio.stratio.com/nexus/content/sites/paas/valkiria/${VERSION}/
else
    echo "target file not available, please run 'make compile' first"
fi

