module mall

go 1.24.9

require (
	github.com/go-chi/chi/v5 v5.2.4
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.4
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/lib/pq v1.10.9
	github.com/nats-io/nats.go v1.48.0
	github.com/pressly/goose/v3 v3.26.0
	github.com/stretchr/testify v1.11.1
	go.uber.org/mock v0.6.0
	golang.org/x/sync v0.19.0
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.18.3 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/nats-io/nkeys v0.4.11 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	go.opentelemetry.io/otel v1.39.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.39.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260114163908-3f89685c29c3 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260114163908-3f89685c29c3 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

tool (
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	google.golang.org/protobuf/cmd/protoc-gen-go
)
