#!/bin/bash

. bin/commons.sh

if [ -d "$GOPATH" ]; then
    if [ ! -d "$SUREFIRE_REPORTS_PATH" ]; then
        mkdir -p $SUREFIRE_REPORTS_PATH
    fi
    cd $GOPATH/src/github.com/Stratio/valkiria
    $GOPATH/bin/godep go test -v -coverprofile=coverageDbus.out ./dbus | $GOPATH/bin/go-junit-report > $SUREFIRE_REPORTS_PATH/reportDbus.xml
    $GOPATH/bin/godep go test -v -coverprofile=coverageProc.out ./proc | $GOPATH/bin/go-junit-report > $SUREFIRE_REPORTS_PATH/reportProc.xml
    $GOPATH/bin/godep go test -v -coverprofile=coveragePluginMesos.out ./plugin/mesos | $GOPATH/bin/go-junit-report > $SUREFIRE_REPORTS_PATH/reportPluginMesos.xml
    $GOPATH/bin/godep go test -v -coverprofile=coverageManager.out ./manager | $GOPATH/bin/go-junit-report > $SUREFIRE_REPORTS_PATH/reportManager.xml
    $GOPATH/bin/godep go test -v -coverprofile=coverageWorkers.out ./workers | $GOPATH/bin/go-junit-report > $SUREFIRE_REPORTS_PATH/reportWorkers.xml
#    $GOPATH/bin/godep go test -v -coverprofile=coverageRoutes.out ./routes | $GOPATH/bin/go-junit-report > $SUREFIRE_REPORTS_PATH/reportRoutes.xml
#    $GOPATH/bin/godep go test -v -coverprofile=coverageValkiria.out ./valkiria | $GOPATH/bin/go-junit-report > $SUREFIRE_REPORTS_PATH/reportValkiria.xml

    # Create combined coverage file
    echo "mode: set" > coverageAll.out
    cat coverageDbus.out | grep -v "mode: set" >> coverageAll.out
    cat coverageProc.out | grep -v "mode: set" >> coverageAll.out
    cat coveragePluginMesos.out | grep -v "mode: set" >> coverageAll.out
    cat coverageManager.out | grep -v "mode: set" >> coverageAll.out
    cat coverageWorkers.out | grep -v "mode: set" >> coverageAll.out
#    cat coverageRoutes.out | grep -v "mode: set" >> coverageAll.out
#    cat coverageValkiria.out | grep -v "mode: set" >> coverageAll.out
else
    echo "target file not available, please run 'make compile' first"
fi
