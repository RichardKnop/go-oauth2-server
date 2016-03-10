# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Install pg client to init db
RUN apt-get update
RUN apt-get install -y postgresql-client-9.4

ADD . /go/src/github.com/RichardKnop/go-oauth2-server

ENV GO15VENDOREXPERIMENT 1
WORKDIR /go/src/github.com/RichardKnop/go-oauth2-server

# Build the go-oauth2-server command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install github.com/RichardKnop/go-oauth2-server

# Copy the docker-entrypoint.sh script and use it as entrypoint
COPY ./docker-compose/docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]

# Document that the service listens on port 8080.
EXPOSE 8080
