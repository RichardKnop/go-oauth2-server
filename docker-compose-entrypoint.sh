#!/bin/bash
set -e

if [ "$1" = 'go-oauth2-server' ] && [ "$2" = 'runserver' ]; then
  curl -L http://etcd:2379/v2/keys/config/go_oauth2_server.json -XPUT -d value='{
    "Database": {
      "Type": "postgres",
      "Host": "postgres",
      "Port": 5432,
      "User": "example_api",
      "Password": "",
      "DatabaseName": "example_api",
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

  $1 migrate
  $1 loaddata oauth/fixtures/scopes.yml
  $1 loaddata oauth/fixtures/test_clients.yml
  $1 loaddata accounts/fixtures/roles.yml
fi

exec "$@"
