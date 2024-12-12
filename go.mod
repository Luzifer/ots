module github.com/Luzifer/ots

go 1.23

toolchain go1.23.0

replace github.com/Luzifer/ots/pkg/customization => ./pkg/customization

require (
	github.com/Luzifer/go_helpers/v2 v2.25.0
	github.com/Luzifer/ots/pkg/customization v0.0.0-20241205141904-ee46f22447b9
	github.com/Luzifer/rconfig/v2 v2.5.2
	github.com/Masterminds/sprig/v3 v3.3.0
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/gorilla/mux v1.8.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.20.5
	github.com/redis/go-redis/v9 v9.7.0
	github.com/sirupsen/logrus v1.9.3
)

require (
	dario.cat/mergo v1.0.1 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.3.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.61.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/spf13/cast v1.7.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	google.golang.org/protobuf v1.35.2 // indirect
	gopkg.in/validator.v2 v2.0.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
