-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA search;

SET
  SEARCH_PATH TO search,
  public;

CREATE TABLE customers_cache (
  id text NOT NULL,
  name text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW (),
  updated_at timestamptz NOT NULL DEFAULT NOW (),
  PRIMARY KEY (id)
);

CREATE TRIGGER created_at_customers_trgr BEFORE
UPDATE ON customers_cache FOR EACH ROW EXECUTE PROCEDURE created_at_trigger ();

CREATE TRIGGER updated_at_customers_trgr BEFORE
UPDATE ON customers_cache FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger ();

CREATE TABLE stores_cache (
  id text NOT NULL,
  name text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW (),
  updated_at timestamptz NOT NULL DEFAULT NOW (),
  PRIMARY KEY (id)
);

CREATE TRIGGER created_at_stores_trgr BEFORE
UPDATE ON stores_cache FOR EACH ROW EXECUTE PROCEDURE created_at_trigger ();

CREATE TRIGGER updated_at_stores_trgr BEFORE
UPDATE ON stores_cache FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger ();

CREATE TABLE products_cache (
  id text NOT NULL,
  store_id text NOT NULL,
  name text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW (),
  updated_at timestamptz NOT NULL DEFAULT NOW (),
  PRIMARY KEY (id)
);

CREATE TRIGGER created_at_products_trgr BEFORE
UPDATE ON products_cache FOR EACH ROW EXECUTE PROCEDURE created_at_trigger ();

CREATE TRIGGER updated_at_products_trgr BEFORE
UPDATE ON products_cache FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger ();

CREATE TABLE orders (
  order_id text NOT NULL,
  customer_id text NOT NULL,
  customer_name text NOT NULL,
  items bytea NOT NULL,
  status text NOT NULL,
  product_ids text ARRAY NOT NULL,
  store_ids text ARRAY NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW (),
  updated_at timestamptz NOT NULL DEFAULT NOW (),
  PRIMARY KEY (order_id)
);

CREATE TRIGGER created_at_orders_trgr BEFORE
UPDATE ON orders FOR EACH ROW EXECUTE PROCEDURE created_at_trigger ();

CREATE TRIGGER updated_at_orders_trgr BEFORE
UPDATE ON orders FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger ();

CREATE TABLE inbox (
  id text NOT NULL,
  name text NOT NULL,
  subject text NOT NULL,
  data bytea NOT NULL,
  received_at timestamptz NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE outbox (
  id text NOT NULL,
  name text NOT NULL,
  subject text NOT NULL,
  data bytea NOT NULL,
  published_at timestamptz,
  PRIMARY KEY (id)
);

CREATE INDEX search_unpublished_idx ON outbox (published_at)
WHERE
  published_at IS NULL;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS search CASCADE;

-- +goose StatementEnd
