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
