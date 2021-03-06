#!/bin/bash

if ! [ "`which docker-compose`" ]; then
  echo 'Error: docker-compose must be installed' >&2
  exit 1
fi

sudo docker-compose exec geth geth --preload=/root/init.js attach
