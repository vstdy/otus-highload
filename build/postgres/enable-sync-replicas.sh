#!/bin/sh
set -e

cat <<-EOM >> "$PGDATA/postgresql.conf"
synchronous_commit = on
synchronous_standby_names = 'FIRST 1 ($POSTGRES_REPLICA_NAME_1, $POSTGRES_REPLICA_NAME_2)'
EOM
