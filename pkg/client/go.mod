module github.com/Luzifer/ots/pkg/client

go 1.25.0

toolchain go1.26.6

replace github.com/Luzifer/ots/pkg/customization => ../customization

require (
	github.com/Luzifer/go-openssl/v4 v4.2.5
	github.com/Luzifer/ots/pkg/customization v0.0.0-20260817110948-81fc004c7ad4
	github.com/ryanuber/go-glob v1.0.0
	github.com/sirupsen/logrus v1.10.0
	github.com/stretchr/testify v1.12.1
)

require (
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/crypto v0.55.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
