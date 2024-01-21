# syntax=docker/dockerfile:1

FROM golang:latest

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY * ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /service

EXPOSE 3000

LABEL org.opencontainers.image.source=https://github.com/notify-fhict/microservices/user

# Run
CMD ["/service"]