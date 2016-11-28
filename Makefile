BINARY = valkiria
BINARYDEBUG = valkiriaDebug
LOGLEVEL = DEBUG
ALLPACKAGE = ./...
PACKAGE = ./valkiria ./routes ./proc ./dbus
VERSION = 0.0.1-SNAPSHOT
# BUILD = `git rev-parse HEAD`
LDFLAGS = -X main.Version=${VERSION} -X main.Build=${BUILD} -X main.LOGLEVEL=${LOGLEVEL}

compile:
	godep go install -v -ldflags "-w $(LDFLAGS)" ${ALLPACKAGE}

test:
	godep go test ${PACKAGE}

package:
	bin/package.sh

docker:
	bin/docker.sh

deploy:
	bin/deploy.sh

env:
	bin/env.sh

buildDebug:
	godep go build -gcflags "-N -l $(LDFLAGS)" ${ALLPACKAGE}

save:
	godep save ${ALLPACKAGE}
