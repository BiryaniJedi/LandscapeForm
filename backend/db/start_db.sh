#!/bin/bash
if [[ "$1" == "test" ]]; then
  echo "test database"
  export $(grep -v '^#' ../.env.testing | xargs) && psql "$DATABASE_URL"
else
  echo "dev database"
  export $(grep -v '^#' ../.env | xargs) && psql "$DATABASE_URL"
  exit 1
fi
