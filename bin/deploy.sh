#!/bin/bash -e

BASEDIR=`dirname $0`/..
VERSION=`cat $BASEDIR/VERSION`
# Upload normal universe
cd target/bin && curl -sS -u stratio:${NEXUSPASS} --upload-file valkiria-${VERSION}.tar.gz http://sodio.stratio.com/nexus/content/sites/paas/valkiria/${VERSION}/
