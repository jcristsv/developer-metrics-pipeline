FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
# O comando tidy gera o go.sum dentro do container antes do build
RUN go mod tidy
RUN go build -o aggregator ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/aggregator .
EXPOSE 8080
CMD ["./aggregator"]
