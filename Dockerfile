FROM golang:1.20-alpine

LABEL maintainer="itelman"
WORKDIR /
COPY . .

# Download Go modules
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY . .

# Build
RUN CGO_ENABLED=1 GOOS=linux go build -o forum
EXPOSE 8080
CMD ["./forum"]
