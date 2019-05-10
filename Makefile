VER_BOOTSTRAP=4.3.1
VER_BOOTSWATCH=4.3.1
VER_FONTAWESOME=5.8.2
VER_GIBBERISH_AES=1.0.0
VER_JQUERY=3.4.1
VER_POPPER=1.15.0
VER_VUE=2.6.10


default: generate

generate: l10n download_libs
	docker run --rm -ti -v $(CURDIR):$(CURDIR) -w $(CURDIR)/src node:10-alpine \
		sh -exc "npm ci && npm run build && rm -rf node_modules && chown -R $(shell id -u) ../frontend"
	go generate

l10n:
	cd frontend/locale && goi18n *

publish:
	curl -sSLo golang.sh https://raw.githubusercontent.com/Luzifer/github-publish/master/golang.sh
	bash golang.sh

# -- Download / refresh external libraries --

clean_libs:
	rm -rf frontend/css frontend/webfonts frontend/js

download_libs: clean_libs fontawesome libs_js libs_css

fontawesome:
	curl -sSfL https://github.com/FortAwesome/Font-Awesome/archive/$(VER_FONTAWESOME).tar.gz | \
		tar -vC frontend -xz --strip-components=1 --wildcards --exclude='*/js-packages' '*/css' '*/webfonts'

libs_css:
	mkdir -p frontend/css
	curl -sSfLo frontend/css/bundle.css "https://cdn.jsdelivr.net/combine/npm/bootstrap@$(VER_BOOTSTRAP)/dist/css/bootstrap.min.css,npm/bootswatch@$(VER_BOOTSWATCH)/dist/flatly/bootstrap.min.css"

libs_js:
	mkdir -p frontend/js
	curl -sSfLo frontend/js/bundle.js "https://cdn.jsdelivr.net/combine/npm/jquery@$(VER_JQUERY),npm/popper.js@$(VER_POPPER),npm/bootstrap@$(VER_BOOTSTRAP)/dist/js/bootstrap.min.js,npm/gibberish-aes@$(VER_GIBBERISH_AES)/dist/gibberish-aes-$(VER_GIBBERISH_AES).min.js,npm/vue@$(VER_VUE)"
