# syntax=docker/dockerfile:1

FROM golang:1.21.4-alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.* ./
RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . ./

# Build

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8888

# Run
ENTRYPOINT CompileDaemon -build="go build ./cmd/api/main.go" -command="./main" -graceful-kill="true"
