VER_FONTAWESOME=5.14.0


default: generate download_libs

generate:
	docker run --rm -ti -v $(CURDIR):$(CURDIR) -w $(CURDIR)/src node:14-alpine \
		sh -exc "npx npm@lts ci && npx npm@lts run build && chown -R $(shell id -u) ../frontend node_modules"

publish: download_libs
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
