#!/bin/sh
set -e

while ! nc -w 1 -zv pg-master 5432; do sleep 1; done
./otus-project migrate up
./otus-project generate --file_path ./people.csv

exec "$@"
