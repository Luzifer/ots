VER_BOOTSTRAP=4.0.0
VER_BOOTSWATCH=4.1.1
VER_FONTAWESOME=5.0.12
VER_GIBBERISH_AES=1.0.0
VER_JQUERY=3.3.1
VER_POPPER=1.14.1


default: generate

generate: l10n download_libs
	docker run --rm -ti -v $(CURDIR):$(CURDIR) -w $(CURDIR) node \
		bash -c "npm install -g coffeescript && coffee -c frontend/application.coffee && chown -R $(shell id -u) frontend"
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
		tar -C frontend -xz --strip-components=2 --wildcards '*/web-fonts-with-css/css' '*/web-fonts-with-css/webfonts'

libs_css:
	mkdir -p frontend/css
	curl -sSfLo frontend/css/bootstrap.min.css "https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/$(VER_BOOTSTRAP)/css/bootstrap.min.css"
	curl -sSfLo frontend/css/bootswatch.min.css "https://cdnjs.cloudflare.com/ajax/libs/bootswatch/$(VER_BOOTSWATCH)/flatly/bootstrap.min.css"

libs_js:
	mkdir -p frontend/js
	curl -sSfLo frontend/js/jquery.min.js "https://cdnjs.cloudflare.com/ajax/libs/jquery/$(VER_JQUERY)/jquery.min.js"
	curl -sSfLo frontend/js/bootstrap.min.js "https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/$(VER_BOOTSTRAP)/js/bootstrap.min.js"
	curl -sSfLo frontend/js/popper.min.js "https://cdnjs.cloudflare.com/ajax/libs/popper.js/$(VER_POPPER)/popper.min.js"
	curl -sSfLo frontend/js/gibberish-aes.min.js "https://cdnjs.cloudflare.com/ajax/libs/gibberish-aes/$(VER_GIBBERISH_AES)/gibberish-aes.min.js"
