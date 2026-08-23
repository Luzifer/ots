module github.com/Luzifer/ots/cmd/ots-cli

go 1.25.0

toolchain go1.27.0

replace (
	github.com/Luzifer/ots/pkg/client => ../../pkg/client
	github.com/Luzifer/ots/pkg/customization => ../../pkg/customization
)

require (
	github.com/Luzifer/ots/pkg/client v0.0.0-20260817110948-81fc004c7ad4
	github.com/sirupsen/logrus v1.10.1
	github.com/spf13/cobra v1.10.2
	github.com/stretchr/testify v1.12.1
)

require (
	github.com/Luzifer/go-openssl/v4 v4.2.5 // indirect
	github.com/Luzifer/ots/pkg/customization v0.0.0-20260817110948-81fc004c7ad4 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/crypto v0.55.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
