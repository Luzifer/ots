module github.com/Luzifer/ots/pkg/client

go 1.23.0

toolchain go1.25.0

replace github.com/Luzifer/ots/pkg/customization => ../customization

require (
	github.com/Luzifer/go-openssl/v4 v4.2.4
	github.com/Luzifer/ots/pkg/customization v0.0.0-20250501151834-283ffa548fa8
	github.com/ryanuber/go-glob v1.0.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
