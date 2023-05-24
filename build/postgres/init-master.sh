#!/bin/sh -e

psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
-v ON_ERROR_STOP=1 -v pg_replica_user="$POSTGRES_REPLICATION_USER" <<-EOSQL
  CREATE ROLE :pg_replica_user WITH LOGIN REPLICATION PASSWORD 'password';
EOSQL

cat <<-EOF >> "$PGDATA/pg_hba.conf"

# New replication connection
host  replication  $POSTGRES_REPLICATION_USER  0.0.0.0/0  trust
EOF

cp /tmp/postgresql.conf /var/lib/postgresql/data/postgresql.conf
