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
	rm -Rfi coverage*.out

install:
	go get github.com/tools/godep
	godep go install -v -ldflags "-w" ./...

cover:
	godep go test -v -coverprofile=coverageDbus.out ./dbus
	godep go test -v -coverprofile=coverageProc.out ./proc
	godep go tool cover -html=coverageDbus.out
	godep go tool cover -html=coverageProc.out

fmt:
	godep go fmt ./proc
	godep go fmt ./dbus

save:
	godep save ./...