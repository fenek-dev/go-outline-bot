# syntax=docker/dockerfile:1

FROM golang:1.22-alpine

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh curl

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@master

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

RUN go build -o /build/app ./cmd/main.go

ENTRYPOINT ["/build/app", "--config=/app/.env"]