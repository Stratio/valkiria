ALLPACKAGE = ./...
PACKAGE = ./valkiria ./routes ./proc ./dbus
LOGLEVEL = DEBUG
VERSION = 0.0.1-SNAPSHOT
BUILD = master
LDFLAGS = -X main.Version=${VERSION} -X main.Build=${BUILD} -X main.LOGLEVEL=${LOGLEVEL}

compile:
	bin/compile.sh

test:
	godep go test ${PACKAGE}

package:
	bin/package.sh

docker:
	bin/docker.sh

deploy:
	bin/deploy.sh
