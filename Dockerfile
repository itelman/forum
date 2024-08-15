FROM golang:1.20.1-alpine3.16 AS build
RUN apk --no-cache add ca-certificates git
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY  . .
RUN apk add build-base
RUN go build -o forum cmd/main.go
FROM alpine:latest
WORKDIR /
COPY --from=build /app .
EXPOSE 8081
CMD ["./forum"]