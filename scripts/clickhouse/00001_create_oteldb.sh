#!/bin/sh
set -e

clickhouse-client -n <<-EOSQL
  CREATE DATABASE oteldb;

  CREATE USER oteldb_user IDENTIFIED WITH sha256_password BY 'oteldb_pass';
  
  GRANT SELECT, INSERT, CREATE DATABASE, CREATE TABLE, CREATE VIEW ON oteldb.* TO oteldb_user;
EOSQL