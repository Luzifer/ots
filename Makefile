PRODUCT_VERSION := v1.21.4

default: build-local

build-local: export CGO_ENABLED=0
build-local: export SOURCE_DATE_EPOCH=1
build-local: frontend generate-apidocs
	go build \
		-buildmode=pie \
		-buildvcs=false \
		-ldflags "-s -w -buildid= -X main.version=$(PRODUCT_VERSION)" \
		-mod=readonly \
		-trimpath \
		-o ots

generate-apidocs: node_modules
	pnpm redocly \
		--disableGoogleFont true \
		-o frontend/api.html \
		build-docs docs/openapi.yaml

frontend_prod: export NODE_ENV=production
frontend_prod: frontend

frontend: node_modules
	pnpm node ci/build.mjs

frontend_lint: node_modules
	pnpm eslint --fix src

node_modules:
	pnpm install --production=false --frozen-lockfile

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
