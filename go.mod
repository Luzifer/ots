module github.com/Luzifer/ots

go 1.25.7

toolchain go1.26.2

replace github.com/Luzifer/ots/pkg/customization => ./pkg/customization

require (
	github.com/Luzifer/go_helpers/file v0.6.1
	github.com/Luzifer/go_helpers/http v0.12.3
	github.com/Luzifer/ots/pkg/customization v0.0.0-20260407120015-d6c630e9a5ea
	github.com/Luzifer/rconfig/v2 v2.6.1
	github.com/Masterminds/sprig/v3 v3.3.0
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/gorilla/mux v1.8.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.23.2
	github.com/redis/go-redis/v9 v9.18.0
	github.com/sirupsen/logrus v1.9.4
)

require (
	dario.cat/mergo v1.0.2 // indirect
	github.com/Luzifer/go_helpers/accesslogger v0.1.1 // indirect
	github.com/Luzifer/go_helpers/str v0.5.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.67.5 // indirect
	github.com/prometheus/procfs v0.20.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.yaml.in/yaml/v2 v2.4.4 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/validator.v2 v2.0.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
