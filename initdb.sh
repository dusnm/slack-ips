#!/usr/bin/env bash

set -e

DB_FILE="./db/ips.db"

mkdir -p "./db"

sqlite3 "$DB_FILE" <<'SQL'
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT,
    name TEXT,
    bank_account_number TEXT,
    city TEXT,
    ips_string TEXT,
    UNIQUE(username),
    UNIQUE(bank_account_number)
);
SQL

echo "Database created at $DB_FILE"