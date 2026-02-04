-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA notifications;

SET
  SEARCH_PATH TO notifications,
  public;

CREATE TABLE notifications.customers_cache (
  id text NOT NULL,
  name text NOT NULL,
  sms_number text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW (),
  updated_at timestamptz NOT NULL DEFAULT NOW (),
  PRIMARY KEY (id)
);

CREATE TRIGGER created_at_customers_trgr BEFORE
UPDATE ON notifications.customers_cache FOR EACH ROW EXECUTE PROCEDURE created_at_trigger ();

CREATE TRIGGER updated_at_customers_trgr BEFORE
UPDATE ON notifications.customers_cache FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger ();

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS notifications CASCADE;

-- +goose StatementEnd
