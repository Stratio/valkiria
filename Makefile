compile:
	bin/compile.sh

test:
	bin/test.sh

package:
	bin/package.sh

docker:
	bin/docker.sh

deploy:
	bin/deploy.sh

clean:
	rm -Rf target

install:
	go get github.com/tools/godep
	godep go install -v -ldflags "-w" ./...

cover:
	godep go test -v -coverprofile=coverage.out ./dbus
	godep go tool cover -html=coverage.out

fmt:
	godep go fmt ./dbus