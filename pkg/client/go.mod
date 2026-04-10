module github.com/Luzifer/ots/pkg/client

go 1.25.0

toolchain go1.26.2

replace github.com/Luzifer/ots/pkg/customization => ../customization

require (
	github.com/Luzifer/go-openssl/v4 v4.2.4
	github.com/Luzifer/ots/pkg/customization v0.0.0-20260407120015-d6c630e9a5ea
	github.com/ryanuber/go-glob v1.0.0
	github.com/sirupsen/logrus v1.9.4
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
