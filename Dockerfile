FROM golang:1.20-alpine AS base

RUN apk add build-base 
WORKDIR /

RUN go env -w CGO_ENABLED=1

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o forum .

FROM alpine:latest
WORKDIR /
COPY --from=base /forum .
COPY --from=base . .

CMD ["./forum"]