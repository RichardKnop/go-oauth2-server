#!/bin/bash
set -e

if [ "$1" = 'go-oauth2-server' ] && [ "$2" = 'runserver' ]; then
  $1 migrate
  $1 loaddata oauth/fixtures/scopes.yml
  $1 loaddata accounts/fixtures/roles.yml
fi

exec "$@"
