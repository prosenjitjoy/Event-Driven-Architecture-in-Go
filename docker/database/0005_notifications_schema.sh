#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "malldb" <<-EOSQL
  CREATE SCHEMA notifications;

  CREATE TABLE notifications.customers_cache(
    id text NOT NULL,
    name text NOT NULL,
    sms_number text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY(id)
  );

  CREATE TRIGGER created_at_customers_trgr BEFORE UPDATE ON notifications.customers_cache FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();

  CREATE TRIGGER updated_at_customers_trgr BEFORE UPDATE ON notifications.customers_cache FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

  GRANT USAGE ON SCHEMA notifications TO malldb_user;

  GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA notifications TO malldb_user;
EOSQL