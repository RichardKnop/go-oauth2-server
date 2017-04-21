#!/bin/bash
set -e

last=$({
  shift $(($#-1))
  echo $1
})

if [ "$last" = 'runserver' ]; then
  $1 migrate
  $1 loaddata oauth/fixtures/scopes.yml
  $1 loaddata oauth/fixtures/roles.yml
fi

exec "$@"
