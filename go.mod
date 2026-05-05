module github.com/AltScore/gothic/v2

go 1.25.6

require (
	github.com/AltScore/money v1.15.1
	github.com/PaesslerAG/jsonpath v0.1.1
	github.com/go-playground/validator/v10 v10.30.2
	github.com/google/uuid v1.6.0
	github.com/labstack/echo/v4 v4.15.2
	github.com/looplab/eventhorizon v0.16.0
	github.com/nsf/jsondiff v0.0.0-20260207060731-8e8d90c4c0ac
	github.com/spf13/viper v1.19.0
	github.com/stretchr/testify v1.11.1
	github.com/subosito/gotenv v1.6.0
	github.com/totemcaf/gollections v0.19.0
	go.mongodb.org/mongo-driver/v2 v2.5.0
	go.uber.org/zap v1.27.0
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56
	golang.org/x/sync v0.20.0
	google.golang.org/genproto v0.0.0-20260316180232-0b37fe3546d5
	google.golang.org/grpc v1.79.3
	gopkg.in/yaml.v3 v3.0.1
)

// Forced versions to fix vulnerability
require golang.org/x/net v0.53.0 // indirect

require (
	github.com/PaesslerAG/gval v1.2.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.13 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/klauspost/compress v1.18.3 // indirect
	github.com/labstack/gommon v0.5.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/magiconair/properties v1.8.10 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.22 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.2.0 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.mongodb.org/mongo-driver v1.12.1 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.42.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.50.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260319201613-d00831a3d3e7 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

// use our fork of eventhorizon until out changes were merged
replace github.com/looplab/eventhorizon => github.com/AltScore/eventhorizon v0.17.0-rc.1.0.20260330195351-1c59b3edc3c1
