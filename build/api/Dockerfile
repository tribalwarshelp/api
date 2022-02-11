FROM golang:1.17.7-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .
ARG VERSION="0.0.0"
RUN apk --no-cache add musl-dev gcc build-base
RUN go install github.com/99designs/gqlgen@v0.14.0
RUN go generate ./...
RUN go build -ldflags="-X 'main.Version=$VERSION'" -o twhelpapi ./cmd/api

######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/twhelpapi .

ENV APP_MODE=production
ENV GIN_MODE=release
EXPOSE 8080

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.2.1/wait ./wait
RUN chmod +x ./wait

CMD ./wait && ./twhelpapi
