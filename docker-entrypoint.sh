#!/bin/bash

set -ex

executable="go-oauth2-server"
cmd="$@"

if [ "$1" = 'runserver' ] || [ "$1" = 'loaddata' ] || [ "$3" = 'runserver' ] || [ "$3" = 'loaddata' ]; then
  extra_args=""
  if [ "$3" = 'runserver' ] || [ "$3" = 'loaddata' ]; then
    extra_args="$1 $2"
  fi

  until $executable $extra_args migrate; do
    >&2 echo "Postgres is unavailable - sleeping"
    sleep 1
  done

  $executable $extra_args loaddata oauth/fixtures/scopes.yml
  $executable $extra_args loaddata oauth/fixtures/roles.yml

  if [[ -z "${FIXTURES}" ]]; then
    echo $FIXTURES | base64 -d > /tmp/fixtures.yml
    $executable $extra_args loaddata /tmp/fixtures.yml
  fi
fi

>&2 echo "Postgres is up - executing command: $cmd"
exec $executable $cmd
