compile:
	bin/compile.sh
	bin/cover.sh

change-version:
	echo "Modifying version to: $(version)"
	echo $(version) > VERSION
	
test:
	bin/test.sh

package:
	bin/package.sh

docker:
	bin/docker.sh

deploy:
	bin/deploy.sh

code-quality:
	bin/codeQuality.sh

clean:
	rm -Rf target

install:
	go get github.com/tools/godep
	godep go install -v -ldflags "-w" ./...

cover:
	bin/cover.sh

fmt:
	godep go fmt ./proc
	godep go fmt ./dbus

save:
	godep save ./...
