ARG tag

# Build Stage
FROM golang:alpine AS builder
WORKDIR /mall
COPY go.* ./
RUN go mod download
COPY . ./
RUN go build -ldflags="-s -w" -v -o baskets-service ./baskets/cmd/service

# Run Stage
FROM gcr.io/distroless/static-debian13:${tag} AS runtime
WORKDIR /app
COPY --from=builder /mall/scripts/wait-for.sh .
COPY --from=builder /mall/baskets-service .
CMD ["/app/baskets-service"]