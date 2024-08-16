FROM golang:1.20-alpine AS base

RUN apk add build-base 
WORKDIR /

RUN go env -w CGO_ENABLED=1

COPY . .
RUN go mod download
RUN go build -o forum .

FROM alpine:latest
WORKDIR /
COPY --from=base . .

EXPOSE 8080

CMD ["./forum"]