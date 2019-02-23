#!/bin/bash

if ! [ "`which docker-compose`" ]; then
  echo 'Error: docker-compose must be installed' >&2
  exit 1
fi

docker-compose exec db psql -Upostgres