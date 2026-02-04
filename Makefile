tools:
	go install tool
	go install go.uber.org/mock/mockgen@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

generate:
	go generate ./...

goose:
	# cd migrations
	# goose create -s <migration_name> sql

integration:
	go test -cover ./... -tags integration