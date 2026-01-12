tools:
	go install tool

generate:
	go generate ./...

postgres:
	podman run --name postgres --hostname postgres -e POSTGRES_PASSWORD=postgres -v ./docker/database:/docker-entrypoint-initdb.d -p 5432:5432 -d postgres:16-alpine