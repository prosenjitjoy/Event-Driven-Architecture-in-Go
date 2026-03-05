-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA notifications;

SET
  SEARCH_PATH TO notifications,
  public;

CREATE TABLE customers_cache (
  id text NOT NULL,
  name text NOT NULL,
  sms_number text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW (),
  updated_at timestamptz NOT NULL DEFAULT NOW (),
  PRIMARY KEY (id)
);

CREATE TRIGGER created_at_customers_trgr BEFORE
UPDATE
  ON customers_cache FOR EACH ROW EXECUTE PROCEDURE created_at_trigger ();

CREATE TRIGGER updated_at_customers_trgr BEFORE
UPDATE
  ON customers_cache FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger ();

CREATE TABLE inbox (
  id text NOT NULL,
  name text NOT NULL,
  subject text NOT NULL,
  data bytea NOT NULL,
  metadata bytea NOT NULL,
  sent_at timestamptz NOT NULL,
  received_at timestamptz NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE outbox (
  id text NOT NULL,
  name text NOT NULL,
  subject text NOT NULL,
  data bytea NOT NULL,
  metadata bytea NOT NULL,
  sent_at timestamptz NOT NULL,
  published_at timestamptz,
  PRIMARY KEY (id)
);

CREATE INDEX notifications_unpublished_idx ON outbox (published_at)
WHERE
  published_at IS NULL;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS notifications CASCADE;

-- +goose StatementEnd