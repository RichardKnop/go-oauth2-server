#!/bin/bash

# waiting til services will start
sleep 20;

# init etcd
curl -L http://etcd:2379/v2/keys/config/go_oauth2_server.json -XPUT -d value='{
    "Database": {
        "Type": "postgres",
        "Host": "pg",
        "Port": 5432,
        "User": "postgres",
        "Password": "",
        "DatabaseName": "go_oauth2_server",
        "MaxIdleConns": 5,
        "MaxOpenConns": 5
    },
    "Oauth": {
        "AccessTokenLifetime": 3600,
        "RefreshTokenLifetime": 1209600,
        "AuthCodeLifetime": 3600
    },
    "Session": {
        "Secret": "test_secret",
        "Path": "/",
        "MaxAge": 604800,
        "HTTPOnly": true
    },
    "IsDevelopment": true
}'

export ETCD_HOST=etcd
export ETCD_PORT=2379


# init pg
createdb -U postgres -h pg go_oauth2_server

# run database migrations
/go/bin/go-oauth2-server migrate

# load fixtures
/go/bin/go-oauth2-server loaddata \
  oauth/fixtures/scopes.yml

/go/bin/go-oauth2-server loaddata \
  oauth/fixtures/test_clients.yml

/go/bin/go-oauth2-server loaddata \
  oauth/fixtures/test_users.yml

/go/bin/go-oauth2-server loaddata \
  oauth/fixtures/test_access_tokens.yml

# finally, run the server
/go/bin/go-oauth2-server runserver
