module github.com/Luzifer/ots

go 1.25.7

toolchain go1.26.5

replace github.com/Luzifer/ots/pkg/customization => ./pkg/customization

replace github.com/Luzifer/ots/pkg/tplfunc => ./pkg/tplfunc

require (
	github.com/Luzifer/go_helpers/file v0.6.2
	github.com/Luzifer/go_helpers/http v0.12.5
	github.com/Luzifer/ots/pkg/customization v0.0.0-20260407120015-d6c630e9a5ea
	github.com/Luzifer/ots/pkg/tplfunc v0.0.0-00010101000000-000000000000
	github.com/Luzifer/rconfig/v2 v2.6.2
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/gorilla/mux v1.8.1
	github.com/prometheus/client_golang v1.24.0
	github.com/redis/go-redis/v9 v9.21.0
	github.com/sirupsen/logrus v1.9.4
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/Luzifer/go_helpers/accesslogger v0.1.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.70.0 // indirect
	github.com/prometheus/procfs v0.21.1 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/sys v0.47.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/validator.v2 v2.0.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
