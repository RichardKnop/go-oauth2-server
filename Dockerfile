# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/RichardKnop/go-oauth2-server

ENV GO15VENDOREXPERIMENT 1
WORKDIR /go/src/github.com/RichardKnop/go-oauth2-server

# Build the go-oauth2-server command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install github.com/RichardKnop/go-oauth2-server

# Run the go-oauth2-server command by default when the container starts.
ENTRYPOINT /go/bin/go-oauth2-server migrate && /go/bin/go-oauth2-server runserver

# Document that the service listens on port 8080.
EXPOSE 8080
