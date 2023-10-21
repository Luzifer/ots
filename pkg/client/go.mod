module github.com/Luzifer/ots/pkg/client

go 1.21.1

replace github.com/Luzifer/ots/pkg/customization => ../customization

require (
	github.com/Luzifer/go-openssl/v4 v4.2.1
	github.com/Luzifer/ots/pkg/customization v0.0.0-00010101000000-000000000000
	github.com/ryanuber/go-glob v1.0.0
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
