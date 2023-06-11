#!/bin/sh -e

psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" --set ON_ERROR_STOP=1 <<-EOSQL
  SELECT citus_set_coordinator_host('pg-coordinator', 5432);
EOSQL
