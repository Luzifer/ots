generate: l10n
	coffee -c frontend/application.coffee
	go generate

l10n:
	cd frontend/locale && goi18n *

publish:
	curl -sSLo golang.sh https://raw.githubusercontent.com/Luzifer/github-publish/master/golang.sh
	bash golang.sh
