#!/bin/bash -e

. bin/compile.sh
. bin/commons.sh

if [ -d "$GOPATH" ]; then
    cd ${GOPATH}/bin && curl -sS -u stratio:${NEXUSPASS} --upload-file valkiria-${VERSION}.tar.gz http://sodio.stratio.com/nexus/content/sites/paas/valkiria/${VERSION}/
else
    echo "target file not available, please run 'make compile' first"
fi

