#!/bin/bash -xv

. bin/commons.sh

ALLPACKAGE="./valkiria ./routes ./proc ./dbus"
LOGLEVEL=INFO
BUILD=$(git rev-parse HEAD)
LDFLAGS="-X main.Version=$VERSION -X main.Build=$BUILD -X main.LOGLEVEL=$LOGLEVEL"

mkdir -p $GOPATH/bin $GOPATH/pkg $GOPATH/src/github.com/Stratio/valkiria
go get -u github.com/tools/godep
go get -u github.com/mattn/goveralls
go get -u github.com/jstemmer/go-junit-report
[ -d $GOPATH/src/github.com/Stratio/valkiria/valkiria ] || ln -s $PWD/valkiria $GOPATH/src/github.com/Stratio/valkiria/valkiria
[ -d $GOPATH/src/github.com/Stratio/valkiria/Godeps ] 	|| ln -s $PWD/Godeps $GOPATH/src/github.com/Stratio/valkiria/Godeps
[ -d $GOPATH/src/github.com/Stratio/valkiria/dbus ] 	|| ln -s $PWD/dbus $GOPATH/src/github.com/Stratio/valkiria/dbus
[ -d $GOPATH/src/github.com/Stratio/valkiria/routes ] 	|| ln -s $PWD/routes $GOPATH/src/github.com/Stratio/valkiria/routes
[ -d $GOPATH/src/github.com/Stratio/valkiria/proc ] 	|| ln -s $PWD/proc $GOPATH/src/github.com/Stratio/valkiria/proc
[ -d $GOPATH/src/github.com/Stratio/valkiria/vendor ] 	|| ln -s $PWD/vendor $GOPATH/src/github.com/Stratio/valkiria/vendor
[ -d $GOPATH/src/github.com/Stratio/valkiria/test ] 	|| ln -s $PWD/test $GOPATH/src/github.com/Stratio/valkiria/test
cd $GOPATH/src/github.com/Stratio/valkiria
$GOPATH/bin/godep go install -v -ldflags "-w $LDFLAGS" $ALLPACKAGE