#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE malldb;

  CREATE USER malldb_user WITH ENCRYPTED PASSWORD 'malldb_pass';

  GRANT CREATE, CONNECT ON DATABASE malldb TO malldb_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "malldb" <<-EOSQL
  GRANT USAGE, CREATE ON SCHEMA public TO malldb_user;
EOSQL