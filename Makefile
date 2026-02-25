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
	podman build -t quay.io/prosenjitjoy/mall-baskets -f ./baskets/Dockerfile --build-arg tag=latest .
	podman build -t quay.io/prosenjitjoy/mall-cosec -f ./cosec/Dockerfile --build-arg tag=latest .
	podman build -t quay.io/prosenjitjoy/mall-customers -f ./customers/Dockerfile --build-arg tag=latest .
	podman build -t quay.io/prosenjitjoy/mall-depot -f ./depot/Dockerfile --build-arg tag=latest .
	podman build -t quay.io/prosenjitjoy/mall-notifications -f ./notifications/Dockerfile --build-arg tag=latest .
	podman build -t quay.io/prosenjitjoy/mall-ordering -f ./ordering/Dockerfile --build-arg tag=latest .
	podman build -t quay.io/prosenjitjoy/mall-payments -f ./payments/Dockerfile --build-arg tag=latest .
	podman build -t quay.io/prosenjitjoy/mall-search -f ./search/Dockerfile --build-arg tag=latest .
	podman build -t quay.io/prosenjitjoy/mall-stores -f ./stores/Dockerfile --build-arg tag=latest .

pg_configmap:
	kubectl create configmap pg-initdb --from-file=./scripts/postgresql

ch_configmap:
	kubectl create configmap ch-initdb --from-file=./scripts/clickhouse

access_postgresql:
	# podman exec -it postgres psql -U <db-user> -d <db-name>
	# kubectl exec -it <pod-name> -- psql -h localhost -U <db-user> -d <db-name>

access_clickhouse:
	# podman exec -it clickhouse clickhouse-client --user <db-user> --password <db-pass>
	# TODO

kubectl_host:
	# kubectl port-forward service/<service-name> <local-port>:<service-port>
