VER_FONTAWESOME=5.14.0


default: generate download_libs

generate:
	docker run --rm -i -v $(CURDIR):$(CURDIR) -w $(CURDIR) node:14-alpine \
		sh -exc "apk add make && make -C src -f ../Makefile generate-inner && chown -R $(shell id -u) frontend src/node_modules"
	docker run --rm -ti -v $(CURDIR):$(CURDIR) -w $(CURDIR) node:14-alpine \
		sh -exc "apk add make && make generate-apidocs && chown -R $(shell id -u) frontend"

generate-apidocs:
	npx @redocly/cli build-docs docs/openapi.yaml --disableGoogleFont true -o /tmp/api.html
	mv /tmp/api.html frontend/

generate-inner:
	npx npm@latest ci
	npx npm@latest run build

publish: download_libs generate-apidocs
	$(MAKE) -C src -f ../Makefile generate-inner
	curl -sSLo golang.sh https://raw.githubusercontent.com/Luzifer/github-publish/master/golang.sh
	bash golang.sh

# -- Download / refresh external libraries --

clean_libs:
	rm -rf \
		frontend/css \
		frontend/js \
		frontend/openssl \
		frontend/webfonts

download_libs: clean_libs
download_libs: fontawesome

fontawesome:
	curl -sSfL https://github.com/FortAwesome/Font-Awesome/archive/$(VER_FONTAWESOME).tar.gz | \
		tar -vC frontend -xz --strip-components=1 --wildcards --exclude='*/js-packages' '*/css' '*/webfonts'
