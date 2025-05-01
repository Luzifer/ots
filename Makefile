default: build-local

build-local: frontend generate-apidocs
	go build \
		-buildmode=pie \
		-ldflags "-s -w -X main.version=$(shell git describe --tags --always || echo dev)" \
		-mod=readonly \
		-trimpath

generate-apidocs:
	npx --yes @redocly/cli build-docs docs/openapi.yaml --disableGoogleFont true -o /tmp/api.html
	mv /tmp/api.html frontend/

frontend_prod: export NODE_ENV=production
frontend_prod: frontend

frontend: node_modules
	corepack yarn@1 node ci/build.mjs

frontend_lint: node_modules
	corepack yarn@1 eslint --fix src

node_modules:
	corepack yarn@1 install --production=false --frozen-lockfile

publish: export NODE_ENV=production
publish: frontend_prod generate-apidocs
	bash ./ci/build.sh

translate:
	cd ci/translate && go run . --write-issue-file

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

.PHONY: node_modules
