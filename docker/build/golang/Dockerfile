# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang:1.12-alpine base image
FROM golang:1.14.1-alpine3.11

ENV GO111MODULE=on

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Add Maintainer Info
LABEL maintainer="Equipe 01 - Hetic"

WORKDIR /go/src/gallery
COPY app /go/src/gallery

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o ./bin/gallery .

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./bin/gallery"]