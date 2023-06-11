#!/bin/sh -e

psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" --set ON_ERROR_STOP=1 <<-EOSQL
  ALTER SYSTEM SET wal_level = logical;
EOSQL
