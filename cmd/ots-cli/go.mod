module github.com/Luzifer/ots/cmd/ots-cli

go 1.21.1

replace github.com/Luzifer/ots/pkg/client => ../../pkg/client

require (
	github.com/Luzifer/ots/pkg/client v0.0.0-20231009165114-f1e303cbe5a7
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.7.0
)

require (
	github.com/Luzifer/go-openssl/v4 v4.2.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
)
