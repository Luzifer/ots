module github.com/Luzifer/ots

go 1.21.1

toolchain go1.21.2

replace github.com/Luzifer/ots/pkg/customization => ./pkg/customization

require (
	github.com/Luzifer/go_helpers/v2 v2.20.0
	github.com/Luzifer/ots/pkg/customization v0.0.0-00010101000000-000000000000
	github.com/Luzifer/rconfig/v2 v2.4.0
	github.com/Masterminds/sprig/v3 v3.2.3
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/pkg/errors v0.9.1
	github.com/redis/go-redis/v9 v9.2.1
	github.com/sirupsen/logrus v1.9.3
)

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	gopkg.in/validator.v2 v2.0.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
