VER_FONTAWESOME:=6.4.0


default: generate download_libs

build-local: download_libs generate-inner generate-apidocs
	go build \
		-buildmode=pie \
		-ldflags "-s -w -X main.version=$(git describe --tags --always || echo dev)" \
		-mod=readonly \
		-trimpath

generate:
	docker run --rm -i -v $(CURDIR):$(CURDIR) -w $(CURDIR) node:18-alpine \
		sh -exc "apk add make && make generate-inner generate-apidocs && chown -R $(shell id -u) frontend node_modules"

generate-apidocs:
	npx --yes @redocly/cli build-docs docs/openapi.yaml --disableGoogleFont true -o /tmp/api.html
	mv /tmp/api.html frontend/

generate-inner:
	npm ci --include=dev
	node ./ci/build.mjs

publish: download_libs generate-inner generate-apidocs
	bash ./ci/build.sh

translate:
	cd ci/translate && go run . --write-issue-file

# -- Download / refresh external libraries --

clean_libs:
	rm -rf \
		frontend/css \
		frontend/js \
		frontend/webfonts

download_libs: clean_libs
download_libs: fontawesome

fontawesome:
	curl -sSfL https://github.com/FortAwesome/Font-Awesome/archive/$(VER_FONTAWESOME).tar.gz | \
		tar -vC frontend -xz --strip-components=1 --wildcards --exclude='*/js-packages' '*/css' '*/webfonts'

# -- Vulnerability scanning --

trivy:
	trivy fs . \
		--dependency-tree \
		--exit-code 1 \
		--format table \
		--ignore-unfixed \
		--quiet \
		--scanners config,license,secret,vuln \
		--severity HIGH,CRITICAL \
		--skip-dirs docs
