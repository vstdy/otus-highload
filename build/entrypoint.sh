#!/bin/sh
set -e

while ! nc -w 1 -zv pg-master 5432; do sleep 1; done
./otus-project migrate up
./otus-project generate users --users_file_path ./people.csv
./otus-project generate friends
./otus-project generate posts --posts_file_path ./posts.txt

exec "$@"
