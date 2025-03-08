FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o ip-sync ./cmd/ip-sync/main.go

FROM alpine:latest
COPY --from=builder /app/ip-sync /ip-sync
CMD ["/ip-sync"]
