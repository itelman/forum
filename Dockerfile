FROM golang:1.20-alpine

LABEL maintainer="itelman"
WORKDIR /
COPY . .
RUN go build -o forum
EXPOSE 8080
CMD ["./forum"]