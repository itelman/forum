# Build Stage
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache build-base gcc musl-dev

WORKDIR /app
RUN go env -w CGO_ENABLED=1

COPY . .
RUN go mod download
RUN go build -o forum ./api

# Final Stage
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app .

EXPOSE 8080
CMD ["./forum"]

