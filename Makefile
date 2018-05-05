generate: l10n
	docker run --rm -ti -v $(CURDIR):$(CURDIR) -w $(CURDIR) node \
		bash -c "npm install -g coffeescript && coffee -c frontend/application.coffee && chown -R $(shell id -u) frontend"
	go generate

l10n:
	cd frontend/locale && goi18n *

publish:
	curl -sSLo golang.sh https://raw.githubusercontent.com/Luzifer/github-publish/master/golang.sh
	bash golang.sh
