#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE orderingdb;

  CREATE USER orderingdb_user WITH ENCRYPTED PASSWORD 'orderingdb_pass';

  GRANT CREATE, CONNECT ON DATABASE orderingdb TO orderingdb_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "orderingdb" <<-EOSQL
  GRANT USAGE, CREATE ON SCHEMA public TO orderingdb_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "orderingdb" <<-EOSQL
  CREATE OR REPLACE FUNCTION created_at_trigger()
  RETURNS TRIGGER AS \$\$
  BEGIN
    NEW.created_at := OLD.created_at;
    RETURN NEW;
  END;
  \$\$ language 'plpgsql';

  CREATE OR REPLACE FUNCTION updated_at_trigger()
  RETURNS TRIGGER AS \$\$
  BEGIN
     IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
        NEW.updated_at = NOW();
        RETURN NEW;
     ELSE
        RETURN OLD;
     END IF;
  END;
  \$\$ language 'plpgsql';
EOSQL