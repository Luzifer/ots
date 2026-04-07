module github.com/Luzifer/ots/cmd/ots-cli

go 1.25.0

toolchain go1.26.1

replace (
	github.com/Luzifer/ots/pkg/client => ../../pkg/client
	github.com/Luzifer/ots/pkg/customization => ../../pkg/customization
)

require (
	github.com/Luzifer/ots/pkg/client v0.0.0-20250501151834-283ffa548fa8
	github.com/sirupsen/logrus v1.9.4
	github.com/spf13/cobra v1.10.2
)

require (
	github.com/Luzifer/go-openssl/v4 v4.2.4 // indirect
	github.com/Luzifer/ots/pkg/customization v0.0.0-20260407120015-d6c630e9a5ea // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
