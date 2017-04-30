#!/bin/sh

set -e

executable="go-oauth2-server"
cmd="$@"

if [ "$1" = 'runserver' ] || [ "$1" = 'loaddata' ]; then
  until $executable migrate; do
    >&2 echo "Postgres is unavailable - sleeping"
    sleep 1
  done

  $executable loaddata oauth/fixtures/scopes.yml
  $executable loaddata oauth/fixtures/roles.yml
fi

>&2 echo "Postgres is up - executing command: $cmd"
exec $executable $cmd
