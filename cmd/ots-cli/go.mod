module github.com/Luzifer/ots/cmd/ots-cli

go 1.23

replace (
	github.com/Luzifer/ots/pkg/client => ../../pkg/client
	github.com/Luzifer/ots/pkg/customization => ../../pkg/customization
)

require (
	github.com/Luzifer/ots/pkg/client v0.0.0-20240922105511-48bf8c9ca4fe
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.8.1
)

require (
	github.com/Luzifer/go-openssl/v4 v4.2.2 // indirect
	github.com/Luzifer/ots/pkg/customization v0.0.0-20240922105511-48bf8c9ca4fe // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
