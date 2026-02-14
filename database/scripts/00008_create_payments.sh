#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE paymentsdb;

  CREATE USER paymentsdb_user WITH ENCRYPTED PASSWORD 'paymentsdb_pass';

  GRANT CREATE, CONNECT ON DATABASE paymentsdb TO paymentsdb_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "paymentsdb" <<-EOSQL
  GRANT USAGE, CREATE ON SCHEMA public TO paymentsdb_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "paymentsdb" <<-EOSQL
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