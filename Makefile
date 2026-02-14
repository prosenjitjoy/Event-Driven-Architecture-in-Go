tools:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install go.uber.org/mock/mockgen@latest 
	go install github.com/bufbuild/buf/cmd/buf@latest 
	go install github.com/pressly/goose/v3/cmd/goose@latest 

generate:
	go generate ./...

goose:
	# cd database/migrations
	# goose create -s <migration_name> sql

integration:
	go test -v -cover ./... -tags integration

build_images:
	podman build -t prosenjitjoy/mall-baskets -f ./baskets/Dockerfile .
	podman build -t prosenjitjoy/mall-cosec -f ./cosec/Dockerfile .
	podman build -t prosenjitjoy/mall-customers -f ./customers/Dockerfile .
	podman build -t prosenjitjoy/mall-depot -f ./depot/Dockerfile .
	podman build -t prosenjitjoy/mall-notifications -f ./notifications/Dockerfile .
	podman build -t prosenjitjoy/mall-ordering -f ./ordering/Dockerfile .
	podman build -t prosenjitjoy/mall-payments -f ./payments/Dockerfile .
	podman build -t prosenjitjoy/mall-search -f ./search/Dockerfile .
	podman build -t prosenjitjoy/mall-stores -f ./stores/Dockerfile .