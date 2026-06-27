module translate

go 1.25.7

toolchain go1.26.4

replace github.com/Luzifer/ots/pkg/tplfunc => ../../pkg/tplfunc

require (
	github.com/Luzifer/ots/pkg/tplfunc v0.0.0-00010101000000-000000000000
	github.com/Luzifer/rconfig/v2 v2.6.2
	github.com/mitchellh/hashstructure/v2 v2.0.2
	github.com/sirupsen/logrus v1.9.4
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/kr/pretty v0.3.1 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/sys v0.33.0 // indirect
	gopkg.in/validator.v2 v2.0.1 // indirect
)
