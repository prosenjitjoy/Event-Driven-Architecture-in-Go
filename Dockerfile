# Build Stage
FROM golang:alpine AS builder
WORKDIR /mall
COPY go.* ./
RUN go mod download
COPY . ./
RUN go build -v -o monolith ./cmd/mall

# Run Stage
FROM alpine:latest AS runtime
WORKDIR /app
COPY --from=builder /mall/wait-for.sh .
RUN chmod +x ./wait-for.sh
COPY --from=builder /mall/monolith .
CMD ["/app/monolith"]